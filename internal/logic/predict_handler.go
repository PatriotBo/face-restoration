package logic

import (
	"encoding/json"
	"face-restoration/internal/constdata"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type PredictResponse struct {
	ImageURL string `json:"imageUrl"`
}

func (m *MiniProgramImpl) Predict(ctx *gin.Context) {
	filename := fmt.Sprintf("temp_%d.png", time.Now().UnixMilli())

	out, err := os.Create("../images/" + filename)
	if err != nil {
		fmt.Printf("predict create file err:%v \n", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, ctx.Request.Body)
	if err != nil {
		fmt.Printf("predict copy file data err:%v \n", err)
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

func (m *MiniProgramImpl) handleUploadFile(ctx *gin.Context) {
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
}
