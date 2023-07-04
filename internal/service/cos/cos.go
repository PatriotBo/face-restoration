package cos

import (
	"context"
	"face-restoration/internal/constdata"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Config cos configuration
type Config struct {
	SecretID  string `yaml:"secretID"`
	SecretKey string `yaml:"secretKey"`
	BucketURL string `yaml:"bucketURL"`
}

// Service cos service
type Service interface {
	PutImage(ctx context.Context, name string) error
}

type serviceImpl struct {
	c *cos.Client
}

// New create a new cos Service
func New(cfg Config) Service {
	u, _ := url.Parse(cfg.BucketURL)
	uri := &cos.BaseURL{BucketURL: u}
	return &serviceImpl{
		c: cos.NewClient(uri, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  cfg.SecretID,
				SecretKey: cfg.SecretKey,
			},
		}),
	}
}

// PutImage upload image to cos
// name is like xxx.png, not full filename
func (s *serviceImpl) PutImage(ctx context.Context, name string) error {
	filename := path.Join(constdata.ImagePath, name)
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open image field err:%v", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Printf("ERROR close file:%s failed err:%v \n", filename, err)
		}
	}()

	rsp, err := s.c.Object.Put(ctx, name, f, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosStorageClass: "STANDARD_IA",
		},
	})
	if err != nil {
		return fmt.Errorf("put object err:%v", err)
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("put object failed code:%d", rsp.StatusCode)
	}
	return nil
}
