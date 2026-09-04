package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 是整个 Leslie Blog 后端的总配置。
//
// 可以把 Config 理解成：
// “程序启动以后，需要使用的所有配置集合”。
//
// 目前我们只有：
//   - Server
//   - MySQL
//
// 后面还会继续增加：
//   - Redis
//   - JWT
//   - MinIO
//   - Meilisearch
//   - Casbin
//   - Logger
//     等。
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
}

// ServerConfig 保存 HTTP Server 相关配置。
type ServerConfig struct {
	// Name 是应用名称。
	Name string `mapstructure:"name"`

	// Host 是 HTTP 服务监听地址。
	Host string `mapstructure:"host"`

	// Port 是 HTTP 服务端口。
	Port int `mapstructure:"port"`

	// Mode 是 Gin 的运行模式。
	//
	// 常见：
	// debug
	// release
	Mode string `mapstructure:"mode"`
}

// MySQLConfig 保存 MySQL 数据库配置。
type MySQLConfig struct {
	// Host 是 MySQL 服务地址。
	Host string `mapstructure:"host"`

	// Port 是 MySQL 服务端口。
	Port int `mapstructure:"port"`

	// Database 是数据库名称。
	Database string `mapstructure:"database"`

	// Username 是数据库用户名。
	Username string `mapstructure:"username"`

	// Password 是数据库密码。
	Password string `mapstructure:"password"`

	// Charset 是数据库字符集。
	Charset string `mapstructure:"charset"`

	// MaxOpenConns 是最大打开连接数。
	MaxOpenConns int `mapstructure:"max_open_conns"`

	// MaxIdleConns 是最大空闲连接数。
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// ConnMaxLifetime 是连接最大生命周期。
	//
	// 这里使用 string，而不是 time.Duration。
	//
	// 例如：
	// "1h"
	// "30m"
	// "10s"
	//
	// 后面在 database 层再使用 time.ParseDuration()
	// 将字符串转换成 time.Duration。
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// Load 从 configs 目录读取配置文件。
func Load() (*Config, error) {

	// 创建一个独立的 Viper 实例。
	v := viper.New()

	// 设置配置文件名称。
	//
	// 不需要写扩展名。
	// Viper 会寻找：
	//
	// config.yaml
	// config.yml
	// config.json
	// 等。
	v.SetConfigName("config")

	// 告诉 Viper 配置文件格式是 YAML。
	v.SetConfigType("yaml")

	// 告诉 Viper 去哪里寻找配置文件。
	v.AddConfigPath("./configs")

	// 读取配置文件。
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	// 创建 Config 对象。
	cfg := &Config{}

	// 将 Viper 读取到的数据转换成我们的 Config struct。
	//
	// YAML：
	//
	// mysql:
	//   host: 127.0.0.1
	//
	// 会被转换成：
	//
	// cfg.MySQL.Host == "127.0.0.1"
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 返回完整配置。
	return cfg, nil
}
