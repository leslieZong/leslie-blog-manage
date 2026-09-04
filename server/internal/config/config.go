package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 是整个 Leslie Blog 后端的总配置。
//
// Config 可以理解成整个应用的“配置中心对象”。
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

// ServerConfig 保存 HTTP Server 相关配置。
type ServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// MySQLConfig 保存 MySQL 数据库配置。
type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// JWTConfig 保存 JWT 认证相关配置。
type JWTConfig struct {

	// Secret 是 JWT 签名密钥。
	//
	// JWT 生成 Token 时需要使用 Secret 进行签名。
	//
	// 后续验证 Token 时，
	// 也需要使用同一个 Secret。
	Secret string `mapstructure:"secret"`

	// Issuer 表示 JWT 的签发者。
	//
	// 例如：
	//
	// leslie-blog
	//
	// 用于标识这个 Token 是由哪个系统生成的。
	Issuer string `mapstructure:"issuer"`

	// ExpireHours 表示 Token 有效时间。
	//
	// 例如：
	//
	// 24
	//
	// 表示 Token 有效期为 24 小时。
	ExpireHours int `mapstructure:"expire_hours"`
}

// Load 从 configs 目录读取配置文件。
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	cfg := &Config{}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return cfg, nil
}
