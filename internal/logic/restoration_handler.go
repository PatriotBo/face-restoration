package logic

import (
	"context"
	"fmt"

	"face-restoration/internal/conf"
	"face-restoration/internal/service/codeformer"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offconfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type messageHandler struct {
	codeFormerService codeformer.Service
	oa                *officialaccount.OfficialAccount
}

func newMessageHandler() *messageHandler {
	return &messageHandler{
		codeFormerService: codeformer.New(),
		oa:                newWechatOfficialAccount(),
	}
}

func newWechatOfficialAccount() *officialaccount.OfficialAccount {
	config := &offconfig.Config{
		AppID:          conf.GetWechatConfig().AppID,
		AppSecret:      conf.GetWechatConfig().AppSecret,
		Token:          conf.GetWechatConfig().Token,
		EncodingAESKey: conf.GetWechatConfig().EncodingAESKey,
		Cache:          cache.NewMemory(), // 使用本地缓存 保存 token

	}
	log.Info(fmt.Sprintf("wechat config:%+v", config))
	//wc := wechat.NewWechat()

	return officialaccount.NewOfficialAccount(config)
}

// HandleMessage receive request from user to restoration an image.
// The image will be sending to codeFormer, it returns a simple reply immediately.
func (h *messageHandler) HandleMessage(ctx *gin.Context) {
	server := h.oa.GetServer(ctx.Request, ctx.Writer)
	server.SkipValidate(false) // 跳过请求合法性检查
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		switch msg.MsgType {
		case message.MsgTypeText:
			return h.handleTextMessage(ctx, msg)
		case message.MsgTypeImage:
			return h.handleImageMessage(ctx, msg)
		default:
			return &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText(fmt.Sprintf("不支持的消息类型 ：%s", msg.MsgType)),
			}
		}
	})
	if err := server.Serve(); err != nil {
		log.Error(fmt.Sprintf("server.Serve failed err:%v", err))
		return
	}
	// 发送回复的消息
	if err := server.Send(); err != nil {
		log.Error(fmt.Sprintf("server.Send failed err:%v", err))
	}
	log.Info(fmt.Sprintf("HandleMessage success %s", server.Token))
}

func (h *messageHandler) handleTextMessage(_ context.Context, msg *message.MixMessage) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(fmt.Sprintf("已收到您发动的消息：%s", msg.Content)),
	}
}

func (h *messageHandler) handleImageMessage(ctx context.Context, msg *message.MixMessage) *message.Reply {
	if err := h.restoration(ctx, msg.GetOpenID(), msg.PicURL); err != nil {
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("照片处理失败，请稍候重试。。。"),
		}
	}
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText("照片处理中，请稍候。。。"),
	}
}

// todo: 将送修请求返回的 ID 保存到DB，待轮询并返回结果
func (h *messageHandler) restoration(ctx context.Context, openID, url string) error {
	log.Info(fmt.Sprintf("restoration openID:%s, url:%s", openID, url))
	id, err := h.codeFormerService.SendPredict(ctx, url)
	if err != nil {
		log.Error(fmt.Sprintf("restoration send predict failed :%v", err))
		return err
	}
	log.Info(fmt.Sprintf("restoration send predict success id:%s", id))
	return nil
}
