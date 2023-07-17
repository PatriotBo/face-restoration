package model

import "time"

type SubscribeType int

const (
	VIPMonth     = 1
	ImagePackage = 2
)

type SubscribeStatus int

const (
	Active    SubscribeStatus = 1
	Expired   SubscribeStatus = 2
	Cancelled SubscribeStatus = 3
)

type PaymentStatus int

const (
	Paid          PaymentStatus = 1
	Pending       PaymentStatus = 2
	UserCancelled PaymentStatus = 3
	Failure       PaymentStatus = 4
)

type Orders struct {
	ID            int64           `gorm:"column:id;primary_key"`
	SubscribeID   string          `gorm:"column:subscribe_id;uniqueIndex:uk_subID"` // 所有订阅请求的唯一标识，相当于自己的订单ID
	UserID        string          `gorm:"column:user_id"`                           // 用户唯一标识 微信 openid
	SubType       SubscribeType   `gorm:"column:sub_type"`
	SubPrice      float64         `gorm:"column:sub_price"`
	SubStartDate  string          `gorm:"column:sub_start_date"` // 订阅开始日期 2023-01-01
	SubEndDate    string          `gorm:"column:sub_end_date"`   // 订阅结束日期 2023-01-31
	SubStatus     SubscribeStatus `gorm:"column:sub_status"`     // 订阅状态
	PaymentMethod int             `gorm:"column:payment_method"` // 支付方式
	PaymentStatus PaymentStatus   `gorm:"column:payment_status"` // 支付状态
	CreateAt      time.Time       `gorm:"column:create_at"`
	UpdateAt      time.Time       `gorm:"column:update_at"`
}

func (o Orders) TableName() string {
	return "t_orders"
}
