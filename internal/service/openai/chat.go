package openai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
)

var (
	chatCache *cache.Cache
)

func init() {
	chatCache = cache.New(24*time.Hour, 25*time.Hour)
}

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
	resp, err := loadCache(messages)
	if err == nil {
		return resp, nil
	}

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

	if err = storeCache(messages, rsp); err != nil {
		fmt.Printf("WARN store cache err:%v \n", err)
	}

	return rsp, nil
}

func storeCache(req any, resp openai.ChatCompletionResponse) error {
	key, err := marshalKey(req)
	if err != nil {
		return err
	}

	chatCache.Set(key, resp, cache.DefaultExpiration)
	return nil
}

func loadCache(req any) (openai.ChatCompletionResponse, error) {
	key, err := marshalKey(req)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	data, exist := chatCache.Get(key)
	if !exist {
		return openai.ChatCompletionResponse{}, errors.New("not founc in cache")
	}

	resp, ok := data.(openai.ChatCompletionResponse)
	if !ok {
		return openai.ChatCompletionResponse{}, errors.New("not valid resp type")
	}
	return resp, nil
}

func marshalKey(req any) (string, error) {
	by, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(by)
	return hex.EncodeToString(hash[:]), nil
}
