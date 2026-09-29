package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Load 加载应用配置。
//
// 配置加载顺序：
//
// 1. 默认值
// 2. YAML 配置文件
// 3. 环境变量
//
// 后面的配置会覆盖前面的配置。
func Load(configFile string) (*Config, error) {
	// 加载本地 .env 文件
	err := godotenv.Load("configs/.env")
	if err != nil {
		return nil, fmt.Errorf(
			"⚠️ 本地 .env 文件未加载（生产环境正常，使用系统环境变量）: %v",
			err,
		)
	}

	v := viper.New()

	// ------------------------------------------------
	// 1. 默认值
	// ------------------------------------------------

	v.SetDefault("app.name", "leslie-blog-server")
	v.SetDefault("app.env", "development")

	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.port", 8080)

	v.SetDefault("mysql.host", "127.0.0.1")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.max_open_conns", 20)
	v.SetDefault("mysql.max_idle_conns", 10)
	v.SetDefault("mysql.conn_max_lifetime", 3600)

	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.expire_hour", 24)

	// ------------------------------------------------
	// 2. 配置文件
	// ------------------------------------------------

	if configFile != "" {

		v.SetConfigFile(configFile)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf(
				"failed to read config: %w",
				err,
			)
		}
	}

	// ------------------------------------------------
	// 3. 环境变量
	// ------------------------------------------------

	v.SetEnvKeyReplacer(
		strings.NewReplacer(
			".",
			"_",
		),
	)

	v.AutomaticEnv()

	// ------------------------------------------------
	// 4. 将配置转换成 Config struct
	// ------------------------------------------------

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal config: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 5. 校验配置
	// ------------------------------------------------

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
