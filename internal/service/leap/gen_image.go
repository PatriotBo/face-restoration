package leap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GenImageRequest for generating images api
type GenImageRequest struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt"`
	Steps          int    `json:"steps"`            // The number of steps to use for the inference.
	Width          int    `json:"width"`            // The width of the image to use for the inference.
	Height         int    `json:"height"`           // The height of the image to use for the inference.
	NumberOfImages int    `json:"number_of_images"` // The number of images to generate for the inference.
	// The higher the prompt strength, the closer the generated image will be to the prompt. Must be between 0 and 30.
	PromptStrength int    `json:"prompt_strength"`
	Seed           int64  `json:"seed"`       // The seed to use for the inference. Must be a positive integer.
	WebhookURL     string `json:"webhookUrl"` // An optional webhook URL that will be called when the model is trained
	// Optionally apply face restoration to the generated images. This will make images of faces look more realistic.
	RestoreFaces  bool `json:"restoreFaces"` // Optionally enhance your prompts automatically to generate better results.
	EnhancePrompt bool `json:"enhance_prompt"`
	// Optionally upscale the generated images. This will make the images look more realistic. The default is x1,
	// which means no upscaling. The maximum is x4.
	UpscaleBy string `json:"upscale_by"`
	Sampler   string `json:"sampler"` // Choose the sampler used for your inference. The default is 'unipc'
}

// GenerateImage format request and call /generate-image api
func (s *serviceImpl) GenerateImage(_ context.Context, modID, prompt string) (*Response, error) {
	req, err := generateImageRequest(modID, prompt, s.token)
	if err != nil {
		return nil, fmt.Errorf("generate request failed:%v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed err:%v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("close resp.Body failed:%v \n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp.Body failed :%v", err)
	}

	response := new(Response)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("unmarshal resp failed:%v", err)
	}
	return response, nil
}

func generateImageRequest(modID, prompt, token string) (*http.Request, error) {
	req := defaultGenerateImageRequest()
	req.Prompt = prompt
	by, _ := json.Marshal(req)
	request, err := http.NewRequest("POST", url(modID), bytes.NewBuffer(by))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", token)
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

// defaultGenerateImageRequest returns a default request for /generate-image api.
// As there are many default values should be set to requests, we defined defaultRequest.
// Some of them may become customized in the future version
func defaultGenerateImageRequest() *GenImageRequest {
	return &GenImageRequest{
		NegativePrompt: `blurry, lowres, ugly, boring, poor lighting, dull,unclear, duplicate, error,low quality,
out of frame, watermark, signature, double faces, two people, multiple people`,
		Steps:          50,
		Width:          512,
		Height:         512,
		NumberOfImages: 2,
		PromptStrength: 10,
		Seed:           time.Now().Unix(),
		WebhookURL:     "",
		RestoreFaces:   false,
		UpscaleBy:      "x2",
		Sampler:        "unipc",
	}
}
