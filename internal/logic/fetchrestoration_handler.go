package logic

import (
	"context"
	"face-restoration/internal/dao"
	"face-restoration/internal/model"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func (h *messageHandler) FetchRestoration(ctx *gin.Context) {
	server := h.oa.GetServer(ctx.Request, ctx.Writer)
	server.SkipValidate(false) // 跳过请求合法性检查
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		// 检查是否为点击自定义菜单栏事件
		if msg.MsgType == message.MsgTypeEvent && msg.Event == message.EventClick {
			// 调用客服消息接口发送消息
			h.handleFetchRequest(ctx, msg.GetOpenID())
		}
		return nil
	})
}

func (h *messageHandler) handleFetchRequest(ctx context.Context, openID string) {
	_, err := h.dao.GetPredictResult(ctx, dao.PredictOption{
		OpenID: openID, Status: int(model.Processing),
	})
	if err != nil {
		fmt.Printf("GetPredictResult faield err:%v \n", err)
		h.sendFailCustomerServiceMessage(openID, "服务异常，请稍候重试")
		return
	}
	//

}

func (h *messageHandler) sendImageCustomerServiceMessage(openID, mediaID string) {
	msg := message.NewCustomerImgMessage(openID, mediaID)
	if err := h.oa.GetCustomerMessageManager().Send(msg); err != nil {
		fmt.Printf("sendImageCustomerServiceMessage faield err:%v \n", err)
	}
}

func (h *messageHandler) sendFailCustomerServiceMessage(openID, text string) {
	msg := message.NewCustomerTextMessage(openID, text)
	if err := h.oa.GetCustomerMessageManager().Send(msg); err != nil {
		fmt.Printf("sendFailCustomerServiceMessage faield err:%v \n", err)
	}
}
