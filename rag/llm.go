package rag

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/hildam/eino-rag/conf"
)

// NewDeepSeekModel 创建 deepseek model
func NewDeepSeekModel(ctx context.Context) (*openai.ChatModel, error) {
	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  conf.GetCfg().DeekSeek.APIKey,
		BaseURL: conf.GetCfg().DeekSeek.BaseURL,
		Model:   conf.AppConfig.Embedder.ModelID,
	})
	if err != nil {
		log.Printf("NewDeepSeekModel failed, init model err: %v\n", err)
		return nil, err
	}
	return model, nil
}
