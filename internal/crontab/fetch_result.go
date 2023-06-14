package crontab

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"net/http"
	"os"
	"time"

	"face-restoration/internal/constdata"
	"face-restoration/internal/dao"
	"face-restoration/internal/model"
	"face-restoration/internal/service/codeformer"

	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type FetchCronImpl struct {
	dao       dao.DBDao
	oa        *officialaccount.OfficialAccount
	cfService codeformer.Service
}

func NewFetchCron(oa *officialaccount.OfficialAccount) *FetchCronImpl {
	return &FetchCronImpl{
		dao:       dao.NewDao(),
		oa:        oa,
		cfService: codeformer.New(),
	}
}

func (f *FetchCronImpl) Start() error {
	c := cron.New()
	if err := c.AddFunc("5/* * * * * *", f.fetch); err != nil {
		panic(err)
	}
	c.Start()
	defer c.Stop()

	return nil
}

func (f *FetchCronImpl) fetch() {
	ctx := context.Background()
	list, err := f.dao.ListProcessingRecords(ctx)
	if err != nil {
		fmt.Printf("list processing records err:%v \n", err)
		return
	}

	for _, r := range list {
		// fetch result from codeForm
		resp, err := f.cfService.GetPrediction(ctx, r.PredictID)
		if err != nil || len(resp.Output) == 0 {
			fmt.Printf("get predition:%s faield err:%v \n", r.PredictID, err)
			continue
		}
		// upload result image to WeChat materials
		localName := formatLocalImageName(r.OpenID, r.ID)
		mediaID, err := f.addMaterial(localName, resp.Output)
		if err != nil {
			fmt.Printf("add material err:%v \n", err)
			continue
		}
		// send result to user
		if err = f.sendImageCustomerServiceMessage(r.OpenID, mediaID); err != nil {
			fmt.Printf("send image customer message err:%v \n", err)
			continue
		}
		// update status,result to db
		upRecord := r
		upRecord.Status = int(model.SendBack)
		upRecord.ResultURL = resp.Output
		upRecord.MediaID = mediaID
		upRecord.UpdateTime = time.Now()
		if err = f.dao.UpdatePredictRecord(ctx, upRecord); err != nil {
			fmt.Printf("update record:%+v failed err:%v \n", upRecord, err)
		}
	}
}

func (f *FetchCronImpl) sendImageCustomerServiceMessage(openID, mediaID string) error {
	msg := message.NewCustomerImgMessage(openID, mediaID)
	if err := f.oa.GetCustomerMessageManager().Send(msg); err != nil {
		fmt.Printf("sendImageCustomerServiceMessage faield err:%v \n", err)
		return err
	}
	return nil
}

func (f *FetchCronImpl) addMaterial(name, url string) (string, error) {
	if err := saveImageLocal(name, url); err != nil {
		fmt.Printf("save image local err:%v \n", err)
		return "", err
	}
	mediaID, _, err := f.oa.GetMaterial().AddMaterial(material.MediaTypeImage, name)
	return mediaID, err
}

func saveImageLocal(name, url string) error {
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

func formatLocalImageName(openID string, id int64) string {
	return fmt.Sprintf("id_%d-openID_%s", id, openID)
}
