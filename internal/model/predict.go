package model

import "time"

// Status predict handle status type
type Status int

const (
	None       Status = 0 // 初始
	Processing Status = 1 // 处理中
	Done       Status = 2 // 完成
	Failed     Status = 3 // 失败
)

// PredictRecord 修复请求记录表
type PredictRecord struct {
	ID         int64     `gorm:"id;primary_key"`
	OpenID     string    `gorm:"openid"`
	ImageURL   string    `gorm:"image_url"`
	Status     int       `gorm:"status"` // 0-初始化 1-处理中 2-完成 3-失败
	PredictID  string    `gorm:"predict_id"`
	ResultURL  string    `gorm:"result_url"`
	MediaID    string    `gorm:"media_id"` // 上传 微信 生成的 mediaID
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
}

// TableName table name
func (PredictRecord) TableName() string {
	return "face_restoration_records"
}
