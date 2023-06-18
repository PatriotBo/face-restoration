package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *MiniProgramImpl) Predict(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request")
		return
	}

	filename := header.Filename
	fmt.Println("File name:", filename)

	out, err := os.Create("../images/" + filename)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to save file")
		return
	}

	ctx.String(http.StatusOK, "File uploaded successfully")
	time.Sleep(2 * time.Second)
}
