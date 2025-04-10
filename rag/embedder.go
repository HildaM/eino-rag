package rag

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/hildam/eino-rag/conf"
)

// NewArkEmbedder 创建 ark embedder
func NewArkEmbedder(ctx context.Context) (*ark.Embedder, error) {
	// 初始化 embedder
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: conf.GetCfg().Embedder.APIKey,
		Model:  conf.GetCfg().Embedder.ModelID,
	})
	if err != nil {
		log.Printf("NewArkEmbedder failed, init embedder err: %v\n", err)
		return nil, err
	}
	return embedder, nil
}
