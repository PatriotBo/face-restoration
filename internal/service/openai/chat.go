package openai

import (
	"context"
	"fmt"
	"time"

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
	defer func(start time.Time) {
		fmt.Printf("INFO ChatCompletion messages:%v \n ", messages)
		fmt.Printf("INFO ChatCompletion cost:%v ", time.Since(start))
	}(time.Now())
	rsp, err := s.cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return rsp, err
	}

	if len(rsp.Choices) == 0 {
		return rsp, fmt.Errorf("chat response choices empty rsp:%v", rsp)
	}
	return rsp, nil
}
