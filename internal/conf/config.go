package conf

import (
	"face-restoration/internal/service/leap"
	"fmt"
	"os"
	"path/filepath"

	"face-restoration/internal/constdata"
	"face-restoration/internal/service/cos"

	"gopkg.in/yaml.v2"
)

var globalConfig *Config

// DBConfig db connection config
type DBConfig struct {
	Host        string `yaml:"host"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Port        int    `yaml:"port"`
	DbName      string `yaml:"dbName"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	Timeout     int    `yaml:"timeout"` // 超时时间 单位：秒
}

// OfficialAccountConfig 公众号相关配置
type OfficialAccountConfig struct {
	AppID          string `yaml:"appID"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
}

// CodeFormerConfig config of code former
type CodeFormerConfig struct {
	Token string `yaml:"token"`
}

// Config 配置信息
type Config struct {
	DB             DBConfig              `yaml:"db"`
	Wechat         OfficialAccountConfig `yaml:"wechat"`
	CodeFormer     CodeFormerConfig      `yaml:"codeFormer"`
	ImageURLPrefix string                `yaml:"imageURLPrefix"`
	Cos            cos.Config            `yaml:"cos"`
	Leap           leap.Config           `yaml:"leap"`
}

func init() {
	filename := filepath.Join(constdata.ConfigPath, "config.yaml")
	by, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("read config err:%v", err))
	}

	c := new(Config)
	if err = yaml.Unmarshal(by, c); err != nil {
		panic(fmt.Errorf("unmarshal config err:%v", err))
	}
	globalConfig = c
}

// GetDSN 获取 db dsn
func GetDSN() string {
	cfg := globalConfig.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.Timeout)
}

// GetWechatConfig ...
func GetWechatConfig() OfficialAccountConfig {
	return globalConfig.Wechat
}

// GetCodeFormerToken ...
func GetCodeFormerToken() string {
	return globalConfig.CodeFormer.Token
}

// GetConfig ...
func GetConfig() *Config {
	return globalConfig
}
