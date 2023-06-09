package conf

import "fmt"

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

// Config 配置信息
type Config struct {
	DB DBConfig `yaml:"db"`
}

func init() {
	defaultConfig()
}

func defaultConfig() {
	globalConfig = &Config{
		DB: DBConfig{
			Host:     "43.156.110.119",
			Port:     3306,
			User:     "admin_new",
			Password: "ZHunbuntu!12",
			DbName:   "building_progress",
		},
	}
}

// GetDSN 获取 db dsn
func GetDSN() string {
	cfg := globalConfig.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.Timeout)
}
