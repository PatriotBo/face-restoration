package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	// 创建一个基本的生产配置的 zap 日志记录器
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log = logger
}

type faceRestorationImpl struct {
	Engine *gin.Engine
}

// NewFaceRestorationImpl create a new impl
func NewFaceRestorationImpl() *faceRestorationImpl {
	e := gin.Default()
	e.GET("/images/:filename", GetImage)
	return &faceRestorationImpl{
		Engine: e,
	}
}
