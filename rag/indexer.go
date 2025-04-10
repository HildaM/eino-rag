package rag

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/hildam/eino-rag/conf"
	"github.com/redis/go-redis/v9"
)

// NewRedisIndexer 创建 redis indexer
func NewRedisIndexer(ctx context.Context, client *redis.Client, embedder *ark.Embedder, prefix string) (*redisInd.Indexer, error) {
	indexer, err := redisInd.NewIndexer(ctx, &redisInd.IndexerConfig{
		Client:           client,
		KeyPrefix:        prefix,
		DocumentToHashes: nil,
		BatchSize:        conf.GetCfg().Indexer.BatchSize,
		Embedding:        embedder,
	})
	if err != nil {
		log.Printf("NewRedisIndexer failed, init indexer err: %v\n", err)
		return nil, err
	}
	return indexer, nil
}

// InitRedisVectorIndex 初始化向量索引
func InitRedisVectorIndex(ctx context.Context, client *redis.Client, indexName, prefix string, dimension int64) error {
	// 检查索引是否存在
	if _, err := client.Do(ctx, "FT.INFO", indexName).Result(); err == nil {
		return nil
	}

	// 创建索引
	indexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"content", "TEXT",
		"vector_content", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", dimension,
		"DISTANCE_METRIC", "COSINE",
	}
	if _, err := client.Do(ctx, indexArgs...).Result(); err != nil {
		log.Printf("InitVectorIndex failed, create index err: %v\n", err)
		return err
	}
	// 二次检测索引
	if _, err := client.Do(ctx, "FT.INFO", indexName).Result(); err != nil {
		log.Printf("InitVectorIndex failed, check index err: %v\n", err)
		return err
	}
	return nil
}
