package leap

import (
	"context"
	"fmt"
	"net/http"
)

// Config leap api configuration
type Config struct {
	Token string `yaml:"token"`
}

// Service leap api service
type Service interface {
	GenerateImage(_ context.Context, modID, prompt string) (*GenImageResponse, error)
	GetImages(_ context.Context, modID, inferenceID string) (*GetImagesResponse, error)
}

type serviceImpl struct {
	client *http.Client
	token  string
}

// New create a new service
func New(cfg Config) Service {
	return &serviceImpl{
		client: http.DefaultClient,
		token:  cfg.Token,
	}
}

var address = "https://api.tryleap.ai/api/v1/images/models/%s/inferences"

func url(modID string) string {
	return fmt.Sprintf(address, modID)
}
