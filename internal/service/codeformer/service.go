package codeformer

import (
	"bytes"
	"context"
	"encoding/json"
	"face-restoration/internal/conf"
	"fmt"
	"io"
	"net/http"

	"face-restoration/internal/constdata"
)

// Service code former service
type Service interface {
	SendPredict(ctx context.Context, image string) (string, error)
	GetPrediction(ctx context.Context, id string) (*GetPredictionResponse, error)
}

type serviceImpl struct {
	client *http.Client
}

// New create a new service
func New() Service {
	return &serviceImpl{
		client: http.DefaultClient,
	}
}

// SendPredictRequest request of sending image to prediction
type SendPredictRequest struct {
	Version string `json:"version"`
	Input   struct {
		Image string `json:"image"`
	} `json:"input"`
}

// SendPredictResponse response of sending image to prediction
type SendPredictResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

const url = "https://api.replicate.com/v1/predictions"

// SendPredict send image to prediction by code former
func (s *serviceImpl) SendPredict(_ context.Context, image string) (string, error) {
	request, err := generateSendPredictRequest(image)
	if err != nil {
		return "", fmt.Errorf("generate request failed :%v", err)
	}
	fmt.Printf("SendPredict request:%+v \n ", request)
	resp, err := s.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("do request faield :%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read resp.body failed :%v", err)
	}

	response := new(SendPredictResponse)
	if err := json.Unmarshal(body, response); err != nil {
		return "", fmt.Errorf("unmarshal response failed :%v", err)
	}

	return response.ID, nil
}

// GetPredictionResponse response of getting predict results
type GetPredictionResponse struct {
	ID     string `json:"id"`
	Output string `json:"output"`
	Status string `json:"status"`
}

// GetPrediction get predict result from code former
func (s *serviceImpl) GetPrediction(_ context.Context, id string) (*GetPredictionResponse, error) {
	request, err := generateGetPredictionRequest(id)
	if err != nil {
		return nil, fmt.Errorf("generate request failed :%v", err)
	}
	resp, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request failed :%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp.Body failed :%v", err)
	}

	response := new(GetPredictionResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("unmarshal resp failed :%v", err)
	}
	return response, nil
}

func generateGetPredictionRequest(id string) (*http.Request, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, id), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", token()) // todo:token 更新
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func generateSendPredictRequest(image string) (*http.Request, error) {
	req := new(SendPredictRequest)
	req.Version = constdata.CodeFormerVersion
	req.Input.Image = image
	by, _ := json.Marshal(req)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(by))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", token()) // todo:token 更新
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func token() string {
	return fmt.Sprintf("Token %s", conf.GetCodeFormerToken())
}
