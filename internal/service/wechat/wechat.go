package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type MiniProgramConfig struct {
	AppID  string `yaml:"app_id"`
	AppKey string `yaml:"app_key"`
}

var code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session"

type Service interface {
	Code2Session(_ context.Context, code string) (*Code2SessionResponse, error)
}

type serviceImpl struct {
	client *http.Client
	config MiniProgramConfig
}

// New create a new service
func New(cfg MiniProgramConfig) Service {
	return &serviceImpl{
		client: http.DefaultClient,
		config: cfg,
	}
}

type Code2SessionRequest struct {
	AppID     string `json:"appid"`
	Secret    string `json:"secret"`
	Code      string `json:"js_code"`
	GrantType string `json:"grant_type"` // 授权类型，此处只需填写 authorization_code
}

type Code2SessionResponse struct {
	ErrCode    int32  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

func (s *serviceImpl) Code2Session(_ context.Context, code string) (*Code2SessionResponse, error) {
	u := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		code2sessionURL,
		url.QueryEscape(s.config.AppID),
		url.QueryEscape(s.config.AppKey),
		url.QueryEscape(code))

	fmt.Printf("INFO Code2Session url:%s \n code:%s \n", u, code)
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp body err:%v", err)
	}

	c2sResp := new(Code2SessionResponse)
	if err = json.Unmarshal(body, c2sResp); err != nil {
		return nil, fmt.Errorf("unmarshal resp body err:%v", err)
	}

	if c2sResp.ErrCode != 0 {
		return nil, fmt.Errorf("request failed code:%d,msg:%s", c2sResp.ErrCode, c2sResp.ErrMsg)
	}
	return c2sResp, nil
}
