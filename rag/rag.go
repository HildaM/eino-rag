package rag

import (
	"context"

	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/hildam/eino-rag/conf"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/document"
	"github.com/redis/go-redis/v9"
)

// Engine rag 引擎
type Engine struct {
	IndexName string // 索引名称
	Prefix    string // 索引前缀
	Dimension int64  // 嵌入维度

	redis    *redis.Client // redis 客户端
	embedder *ark.Embedder // 火山引擎 ark embedding 模型

	FileLoader *file.FileLoader      // 文件加载器
	Splitter   *document.Transformer // 文本分割器
	Retriever  *redisRet.Retriever   // redis 检索器
	Indexer    *redisInd.Indexer     // redis 索引器
	LLM        *openai.ChatModel     // 大模型
}

// NewEngine 创建 rag 引擎
func NewEngine(ctx context.Context, index, prefix string) *Engine {
	return &Engine{
		IndexName: index,
		Prefix:    prefix,
		Dimension: conf.GetCfg().Rag.Dimension,
		redis: redis.NewClient(&redis.Options{
			Addr:     conf.GetCfg().Redis.Addr,
			Password: conf.GetCfg().Redis.Password,
		}),

		// TODO: 待设计初始化函数
		embedder:   nil,
		FileLoader: nil,
		Splitter:   nil,
		Retriever:  nil,
		Indexer:    nil,
		LLM:        nil,
	}
}
