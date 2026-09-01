package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	GitHub GitHubConfig `mapstructure:"github"`
	Upload UploadConfig `mapstructure:"upload"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"` // debug / release
	Port      int    `mapstructure:"port"`
	BaseURL   string `mapstructure:"base_url"`
	UploadDir string `mapstructure:"upload_dir"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpireSec int    `mapstructure:"expire_sec"` // token 有效期（秒）
	Issuer    string `mapstructure:"issuer"`
}

type GitHubConfig struct {
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"` // 可选，提升速率限制
}

type UploadConfig struct {
	MaxSizeMB int `mapstructure:"max_size_mb"`
}

var cfg *Config

// Load 读取配置：先加载 config/app.yaml，再用环境变量覆盖（LESLIE_ 前缀 + 下划线）
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("LESLIE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg = &c
	return cfg, nil
}

// Get 取已加载配置
func Get() *Config {
	return cfg
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "Leslie Blog")
	v.SetDefault("app.mode", "debug")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.base_url", "http://localhost:8080")
	v.SetDefault("app.upload_dir", "./uploads")

	v.SetDefault("mysql.host", "127.0.0.1")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "root")
	v.SetDefault("mysql.password", "")
	v.SetDefault("mysql.database", "leslie_blog")
	v.SetDefault("mysql.charset", "utf8mb4")
	v.SetDefault("mysql.max_idle_conns", 10)
	v.SetDefault("mysql.max_open_conns", 100)
	v.SetDefault("mysql.log_level", "warn")

	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.secret", "change-me-please")
	v.SetDefault("jwt.expire_sec", 7*24*3600)
	v.SetDefault("jwt.issuer", "leslie-blog")

	v.SetDefault("upload.max_size_mb", 20)
}
