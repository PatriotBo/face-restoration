package leap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetImages get inference images by calling /get-images api
func (s *serviceImpl) GetImages(_ context.Context, modID, inferenceID string) (*Response, error) {
	request, err := generateGetImageRequest(modID, inferenceID, s.token)
	if err != nil {
		return nil, fmt.Errorf("generate get image request failed:%v", err)
	}

	resp, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do get image request faield:%v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("close resp.Body failed :%v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp.Body failed:%v", err)
	}

	response := new(Response)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("unmarshal resp failed:%v", err)
	}
	return response, nil
}

func generateGetImageRequest(modID, inferenceID, token string) (*http.Request, error) {
	url := fmt.Sprintf("https://api.tryleap.ai/api/v1/images/models/%s/inferences/%s",
		modID, inferenceID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", token)
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}
