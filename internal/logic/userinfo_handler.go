package logic

type PayStatus int

const (
	None      PayStatus = iota // 非付费用户
	PayVIP                     // VIP 用户
	PayImages                  // 购买图片包用户
)

type UserInfo struct {
	Openid     string    `json:"openid"`
	SessionKey string    `json:"session_key"`
	Nickname   string    `json:"nickname"`
	AvatarUrl  string    `json:"avatarUrl"`
	PayStatus  PayStatus `json:"pay_status"`
}
