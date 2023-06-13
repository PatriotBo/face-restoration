package logic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"face-restoration/internal/constdata"
	"face-restoration/internal/dao"
	"face-restoration/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/material"
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
	if err := server.Serve(); err != nil {
		log.Error(fmt.Sprintf("fetch server.Serve failed err:%v", err))
		return
	}
	// 发送回复的消息
	if err := server.Send(); err != nil {
		log.Error(fmt.Sprintf("fetch server.Send failed err:%v", err))
	}
	log.Info(fmt.Sprintf("FetchRestoration success %s", server.Token))
}

// 处理用户获取修复结果的请求
func (h *messageHandler) handleFetchRequest(ctx context.Context, openID string) {
	records, err := h.dao.GetPredictResult(ctx, dao.PredictOption{
		OpenID: openID, Status: int(model.Processing),
	})
	if err != nil {
		fmt.Printf("GetPredictResult faield err:%v \n", err)
		h.sendFailCustomerServiceMessage(openID, "服务异常，请稍候重试")
		return
	}
	for _, v := range records {
		rsp, err := h.codeFormerService.GetPrediction(ctx, v.PredictID)
		if err != nil {
			fmt.Printf("get prediction failed err:%v \n", err)
			continue
		}

		upRecord := v
		upRecord.Status = int(model.Done)
		if rsp.Status != "success" || len(rsp.Output) == 0 {
			fmt.Printf("get prediction failed status:%s \n", rsp.Status)
			upRecord.Status = int(model.Failed)
		}
		upRecord.ResultURL = rsp.Output

		// 图片保存到本地
		imageName := formatLocalImageName(v.OpenID, v.ID)
		if err := h.saveImageLocal(imageName, rsp.Output); err != nil {
			fmt.Printf("save image local failed err:%v \n", err)
			continue
		}
		// 图片上传微信，生成临时素材
		mediaID, err := h.generateMaterial(imageName)
		if err != nil {
			fmt.Printf("generate material failed err:%v \n", err)
			continue
		}
		upRecord.MediaID = mediaID
		//	保存结果到DB
		if err = h.dao.UpdatePredictRecord(ctx, upRecord); err != nil {
			fmt.Printf("save prediction result err:%v \n", err)
			continue
		}
		// 发送客服消息，将修复后的图片返回给用户
		h.sendImageCustomerServiceMessage(v.OpenID, mediaID)
	}
}

func (h *messageHandler) saveImageLocal(name, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("body close faield err:%v \n", err)
		}
	}()

	filename := fmt.Sprintf("%s/%s.png", constdata.ImagePath, name)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("file close failed err:%v \n", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (h *messageHandler) generateMaterial(name string) (string, error) {
	filename := fmt.Sprintf("%s/%s.png", constdata.ImagePath, name)
	mediaID, _, err := h.oa.GetMaterial().AddMaterial(material.MediaTypeImage, filename)
	return mediaID, err
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

func formatLocalImageName(openID string, id int64) string {
	return fmt.Sprintf("id_%d-openID_%s", id, openID)
}
