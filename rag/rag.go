package rag

import (
	"context"
	"log"

	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/google/uuid"
	"github.com/hildam/eino-rag/conf"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

var systemPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

here's documents searched for you:
==== doc start ====
	  {documents}
==== doc end ====
`

// Engine rag 引擎
type Engine struct {
	IndexName string // 索引名称
	Prefix    string // 索引前缀
	Dimension int64  // 嵌入维度

	redis      *redis.Client        // redis 客户端
	FileLoader *file.FileLoader     // 文件加载器
	Splitter   document.Transformer // 文本分割器
	Retriever  *redisRet.Retriever  // redis 检索器
	Indexer    *redisInd.Indexer    // redis 索引器
	LLM        *openai.ChatModel    // 大模型
}

// NewEngine 创建 rag 引擎
func NewEngine(ctx context.Context, index, prefix string) (*Engine, error) {
	// 初始化 redis
	redisCli := redis.NewClient(&redis.Options{
		Addr:     conf.GetCfg().Redis.Addr,
		Password: conf.GetCfg().Redis.Password,
	})

	// 初始化 embedder
	embedder, err := NewArkEmbedder(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化 fileloader
	fileLoader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,
		Parser:      nil,
	})
	if err != nil {
		log.Printf("NewEngine failed, init fileloader err: %v\n", err)
		return nil, err
	}

	// 初始化 splitter
	spliter, err := NewMarkdownSplitter(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化 retriever
	retriever, err := NewRedisRetriever(ctx, redisCli, embedder, index)
	if err != nil {
		return nil, err
	}

	// 初始化 indexer
	indexer, err := NewRedisIndexer(ctx, redisCli, embedder, prefix)
	if err != nil {
		return nil, err
	}

	// 初始化 llm
	llm, err := NewDeepSeekModel(ctx)
	if err != nil {
		return nil, err
	}

	return &Engine{
		IndexName:  index,
		Prefix:     prefix,
		Dimension:  conf.GetCfg().Rag.Dimension,
		redis:      redisCli,
		FileLoader: fileLoader,
		Splitter:   spliter,
		Retriever:  retriever,
		Indexer:    indexer,
		LLM:        llm,
	}, nil
}

// AddFile 添加文件
func (e *Engine) AddFile(ctx context.Context, filepath string) error {
	// 加载文件
	docs, err := e.FileLoader.Load(ctx, document.Source{
		URI: filepath,
	})
	if err != nil {
		log.Printf("CreateFileIndex failed, load file err: %v\n", err)
		return err
	}
	// 分割文本
	docs, err = e.Splitter.Transform(ctx, docs)
	if err != nil {
		log.Printf("CreateFileIndex failed, split text err: %v\n", err)
		return err
	}

	// 为每个文档生成唯一 id
	for _, d := range docs {
		uuid, _ := uuid.NewUUID()
		d.ID = uuid.String()
	}

	// 初始化向量索引
	if err := InitRedisVectorIndex(ctx, e.redis, e.IndexName, e.Prefix, e.Dimension); err != nil {
		log.Printf("CreateFileIndex failed, init vector index err: %v\n", err)
		return err
	}

	// 存储索引
	if _, err := e.Indexer.Store(ctx, docs); err != nil {
		log.Printf("CreateFileIndex failed, store index err: %v\n", err)
		return err
	}
	return nil
}

// Query 查询
func (e *Engine) Query(ctx context.Context, query string) (*schema.StreamReader[*schema.Message], error) {
	// 检索
	docs, err := e.Retriever.Retrieve(ctx, query)
	if err != nil {
		log.Printf("Query failed, retrieve err: %v\n", err)
		return nil, err
	}
	log.Printf("Query success, docs: %v\n\n", docs)

	// 生成 prompt
	promptTempalte := prompt.FromMessages(schema.FString, []schema.MessagesTemplate{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("question: {content}"),
	}...)
	message, err := promptTempalte.Format(ctx, map[string]any{
		"content":   query,
		"documents": docs,
	})
	if err != nil {
		log.Printf("Query failed, format prompt err: %v\n", err)
		return nil, err
	}

	// 调用 llm
	return e.LLM.Stream(ctx, message)
}
