package logic

import (
	"encoding/json"
	"fmt"

	"face-restoration/internal/constdata"

	"github.com/gin-gonic/gin"
)

type PredictResponse struct {
	ImageURL string `json:"imageUrl"`
}

func (m *MiniProgramImpl) Predict(ctx *gin.Context) {
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		fmt.Printf("predict bad request err:%v \n", err)
		return
	}

	filename := header.Filename
	fmt.Println("File name:", filename)
	if err = ctx.SaveUploadedFile(header, constdata.ImagePath); err != nil {
		fmt.Printf("predict save field err:%v \n", err)
		return
	}

	resp := &PredictResponse{
		ImageURL: fmt.Sprintf("%s/.%s", "", filename),
	}
	body, _ := json.Marshal(resp)
	if _, err = ctx.Writer.Write(body); err != nil {
		fmt.Printf("write response err:%v \n", err)
	}
}
