package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Config openai api configuration
type Config struct {
	AuthToken string `yaml:"authToken"`
	BaseURL   string `yaml:"baseURL"`
}

// Service openai api service
type Service interface {
	ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (
		openai.ChatCompletionResponse, error)
}

type serviceImpl struct {
	cli *openai.Client
}

// New create a new openai api client
func New(cfg Config) Service {
	config := openai.DefaultConfig(cfg.AuthToken)
	config.BaseURL = cfg.BaseURL
	return &serviceImpl{
		cli: openai.NewClientWithConfig(config),
	}
}

// ChatCompletion chat with openai api
func (s *serviceImpl) ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (
	openai.ChatCompletionResponse, error) {
	rsp, err := s.cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4,
		Messages: messages,
	})
	if err != nil {
		return rsp, fmt.Errorf("chat failed:%v", err)
	}

	if len(rsp.Choices) == 0 {
		return rsp, fmt.Errorf("chat response choices empty rsp:%v", rsp)
	}
	return rsp, nil
}
