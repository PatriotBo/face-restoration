package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"face-restoration/internal/conf"
	"face-restoration/internal/constdata"

	"github.com/gin-gonic/gin"
)

type PredictResponse struct {
	ImageURL string `json:"imageUrl"`
}

func (m *MiniProgramImpl) Predict(ctx *gin.Context) {
	filename := fmt.Sprintf("temp_%d.png", time.Now().UnixMilli())
	if err := receiveImage(filename, ctx.Request.Body); err != nil {
		return
	}

	output, err := m.predict(ctx, genImageURL(filename))
	if err != nil {
		return
	}

	resp := &PredictResponse{
		ImageURL: genImageCosURL(output),
	}
	fmt.Printf("cos url:%s \n", genImageCosURL(output))
	body, _ := json.Marshal(resp)
	if _, err := ctx.Writer.Write(body); err != nil {
		fmt.Printf("write response err:%v \n", err)
	}
}

func receiveImage(filename string, file io.Reader) error {
	out, err := os.Create("../images/" + filename)
	if err != nil {
		fmt.Printf("predict create file err:%v \n", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Printf("predict copy file data err:%v \n", err)
		return err
	}
	return nil
}

func (m *MiniProgramImpl) predict(ctx context.Context, url string) (string, error) {
	id, err := m.cfService.SendPredict(ctx, url)
	if err != nil {
		fmt.Printf("send predict err:%v \n", err)
		return "", err
	}
	startf := time.Now()
	output, err := m.fetchPredictResult(ctx, id)
	if err != nil {
		fmt.Printf("fetch faild err:%v \n", err)
		return "", err
	}
	fmt.Printf("fetch result cost:%v \n", time.Since(startf))

	startd := time.Now()
	filename, err := downloadPrediction(id, output)
	if err != nil {
		fmt.Printf("download prediction err:%v \n", err)
		return "", err
	}
	fmt.Printf("download image cost:%v \n", time.Since(startd))

	// 上传到cos
	start := time.Now()
	if err = m.cosService.PutImage(ctx, filename); err != nil {
		fmt.Printf("put image to cos err:%v \n", err)
		return "", err
	}
	fmt.Printf("cos put cost:%v \n", time.Since(start))
	return filename, nil
}

func (m *MiniProgramImpl) fetchPredictResult(ctx context.Context, id string) (string, error) {
	fmt.Printf("fetch predict result id:%s \n \n", id)
	t := time.NewTicker(3 * time.Second)
	count := 20
	for range t.C {
		count--
		if count < 0 {
			return "", errors.New("fetch predict result failed")
		}
		rsp, err := m.cfService.GetPrediction(ctx, id)
		if err != nil {
			return "", err
		}
		if rsp.Status != "succeeded" || len(rsp.Output) == 0 {
			fmt.Printf("prediction not ready status:%s \n", rsp.Status)
			continue
		}
		return rsp.Output, nil
	}
	return "", nil
}

func downloadPrediction(id, url string) (string, error) {
	fmt.Printf("downloadPrediction begin to download images id:%s \n", id)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("body close faield err:%v \n", err)
		}
	}()

	filename := fmt.Sprintf("%s/%s.png", constdata.ImagePath, id)
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("file close failed err:%v \n", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	return fmt.Sprintf("%s.png", id), err
}

func genImageURL(name string) string {
	return fmt.Sprintf("%s/%s", conf.GetConfig().ImageURLPrefix, name)
}

func genImageCosURL(name string) string {
	return fmt.Sprintf("%s/%s", conf.GetConfig().Cos.BucketURL, name)
}
