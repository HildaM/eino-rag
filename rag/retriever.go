package rag

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/redis/go-redis/v9"
)

// NewRedisRetriever 创建 redis retriever
func NewRedisRetriever(ctx context.Context, client *redis.Client, embedder *ark.Embedder,
	index string) (*redisRet.Retriever, error) {

	retriever, err := redisRet.NewRetriever(ctx, &redisRet.RetrieverConfig{
		Client:            client,
		Index:             index,
		Embedding:         embedder,
		VectorField:       "vector_content",
		DistanceThreshold: nil,
		Dialect:           2,
		ReturnFields:      []string{"vector_content", "content"},
		DocumentConverter: nil,
		TopK:              1,
	})
	if err != nil {
		log.Printf("NewRedisRetriever failed, init retriever err: %v\n", err)
		return nil, err
	}
	return retriever, nil
}
