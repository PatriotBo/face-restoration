package logic

import (
	"fmt"
	"path/filepath"

	"face-restoration/internal/constdata"

	"github.com/gin-gonic/gin"
)

// GetImage handler
func GetImage(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filePath := filepath.Join(constdata.ImagePath, filename)

	log.Info(fmt.Sprintf("GetImage filepath:%s", filePath))
	// 设置响应头，使浏览器能够识别为图片文件
	ctx.Header("Content-Type", "image/jpeg")
	// 设置响应头，使浏览器能够识别为附件下载
	ctx.Header("Content-Disposition", "attachment; filename="+filename)

	// 发送图片文件
	ctx.File(filePath)
}
