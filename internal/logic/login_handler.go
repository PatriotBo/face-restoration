package logic

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// LoginRequest ...
type LoginRequest struct {
	Code string `json:"code"`
}

// LoginResponse ...
type LoginResponse struct {
	Openid     string    `json:"openid"`
	SessionKey string    `json:"session_key"`
	PayStatus  PayStatus `json:"pay_status"`
}

// Login api for user login
func (m *MiniProgramImpl) Login(ctx *gin.Context) {
	request := new(LoginRequest)
	if err := ctx.BindJSON(request); err != nil {
		fmt.Printf("ERROR login parse request err:%v \n", err)
		ctx.JSON(400, "bad request")
		return
	}
	fmt.Printf("INFO login request:%+v \n", request)
	resp, err := m.wechatService.Code2Session(ctx, request.Code)
	if err != nil {
		fmt.Printf("ERROR code2session failed:%v \n", err)
		ctx.JSON(500, err.Error())
		return
	}

	response := &LoginResponse{
		Openid:     resp.Openid,
		SessionKey: resp.SessionKey,
		PayStatus:  PayVIP, // todo: 从数据库获取用户账号信息
	}
	fmt.Printf("INFO login resp:%+v \n", response)
	ctx.JSON(200, response)
}
