package logic

import (
	"crypto/tls"
	"face-restoration/internal/constdata"
	"face-restoration/internal/crontab"
	"net/http"

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
	Cron   *crontab.FetchCronImpl
}

// RunService create a new impl
func RunService() {
	e := gin.Default()
	// 静态图片访问
	e.Static("/img", constdata.ImagePath)

	//msgHandler := newMessageHandler()
	//// 微信消息
	//e.POST("/wx", func(ctx *gin.Context) {
	//	msgHandler.HandleMessage(ctx)
	//})

	// 设置HTTPS证书和密钥
	certFile := "../certs/cert.pem"
	keyFile := "../certs/key.pem"

	// 配置TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   e,
		TLSConfig: tlsConfig,
	}
	// 启动HTTPS服务器
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}
