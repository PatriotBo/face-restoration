package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"

	"github.com/gin-gonic/gin"
)

// GenType type of generate images
type GenType int

// WallPaper enum of GenType
const (
	WallPaper GenType = 1
)

// GenerateImageRequest request for generating images
type GenerateImageRequest struct {
	GenType GenType `json:"gen_type"`
	Prompt  string  `json:"prompt"`
}

// GenerateImage generate images handler
func (m *MiniProgramImpl) GenerateImage(ctx *gin.Context) {
	req, err := parseGenerateImageRequest(ctx)
	if err != nil {
		fmt.Printf("ERROR GenerateImage parse req:%v \n", err)
		return
	}
	fmt.Printf("GenerateImage request:%+v \n", req)

	optimizePrompt, err := m.optimizePrompt(ctx, req.Prompt)
	if err != nil {
		fmt.Printf("ERROR GenerateImage optimize prompt:%v", err)
		return
	}
	fmt.Printf("INFO GenerateImage optimize prompt:%s", optimizePrompt)
}

func parseGenerateImageRequest(ctx *gin.Context) (*GenerateImageRequest, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("read req.body failed:%v", err)
	}

	req := new(GenerateImageRequest)
	if err = json.Unmarshal(body, req); err != nil {
		return nil, fmt.Errorf("unmarshal req failed:%v", err)
	}
	return req, nil
}

func (m *MiniProgramImpl) optimizePrompt(ctx context.Context, prompt string) (string, error) {
	transResp, err := m.openAIService.ChatCompletion(ctx, generateTranslateChatMessages(prompt))
	if err != nil {
		return "", fmt.Errorf("translate failed:%v", err)
	}

	transPrompt := transResp.Choices[0].Message.Content
	fmt.Printf("INFO translate prompt:%s \n", transPrompt)

	resp, err := m.openAIService.ChatCompletion(ctx, generatePromptOptimizeChatMessages(transPrompt))
	if err != nil {
		return "", fmt.Errorf("optimize prompt failed:%v", err)
	}
	return resp.Choices[0].Message.Content, nil
}

func generateTranslateChatMessages(prompt string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(translateMessageContentFormat, prompt),
		},
	}
}

func generatePromptOptimizeChatMessages(prompt string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(promptOptimizeMessageContentFormat, prompt),
		},
	}
}

var translateMessageContentFormat = `请将下面的内容翻译为英文："%s"`

var promptOptimizeMessageContentFormat = `
假如你是一个 AI prompt 优化专家，请参考下面的例子，将我后面输入的简单的 prompt 进行优化填充，生成更加详细合理的 prompt，以便 midjourney ，stable diffusion 等 
可以更好的生成图片，可以供你参考的例子如下：
"Mood: Mystical | 
Subject: A mesmerizing labyrinth of crystalline ice caves under a frosty night sky | 
Timing: Midnight | 
Lens: Wide-angle | 
Lighting conditions: Soft, ethereal glow from the moon casting bluish hues on the icy surfaces | 
Style: Fusion of nature and abstract geometry | 
Colors: Sparkling silvers, icy blues, and deep indigos of the night sky | 
Background: A vast, star-studded sky visible through the opening of the ice cave | 
Perspective: Within the cave, looking towards the entrance | 
Focal point: A beautifully formed icicle, its surface reflecting the moonlight | 
Space: Captivating and otherworldly | Pattern/Texture: The intricate, natural patterns of the ice cave walls and the smooth, reflective ice floor | 
Element defining the scale: A solitary snowflake, caught on the tip of the icicle | 
Depth of field: Deep, capturing the enchanting ice cave and the infinite expanse of the night sky | 
Feeling: Intriguing and awe-inspiring | 
Contrast elements: The cold, enduring beauty of the ice cave against the distant, tranquil presence of the starry sky."

需要优化的prompt是：
"%s"
`
