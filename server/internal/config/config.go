package config

import (
	"fmt"
	"strings"
)

// Validate 校验配置是否满足应用启动要求。
func (c *Config) Validate() error {

	// ------------------------------------------------
	// App
	// ------------------------------------------------

	if strings.TrimSpace(c.App.Name) == "" {
		return fmt.Errorf(
			"app.name cannot be empty",
		)
	}

	if strings.TrimSpace(c.App.Env) == "" {
		return fmt.Errorf(
			"app.env cannot be empty",
		)
	}

	// ------------------------------------------------
	// HTTP
	// ------------------------------------------------

	if c.Server.Port <= 0 ||
		c.Server.Port > 65535 {

		return fmt.Errorf(
			"http.port must be between 1 and 65535",
		)
	}

	// ------------------------------------------------
	// MySQL
	// ------------------------------------------------

	if strings.TrimSpace(c.MySQL.Host) == "" {
		return fmt.Errorf(
			"mysql.host cannot be empty",
		)
	}

	if strings.TrimSpace(c.MySQL.Username) == "" {
		return fmt.Errorf(
			"mysql.username cannot be empty",
		)
	}

	if strings.TrimSpace(c.MySQL.Database) == "" {
		return fmt.Errorf(
			"mysql.database cannot be empty",
		)
	}

	// ------------------------------------------------
	// Redis
	// ------------------------------------------------

	if strings.TrimSpace(c.Redis.Addr) == "" {
		return fmt.Errorf(
			"redis.addr cannot be empty",
		)
	}

	// ------------------------------------------------
	// JWT
	// ------------------------------------------------

	if c.App.Env == "production" &&
		strings.TrimSpace(c.JWT.Secret) == "" {

		return fmt.Errorf(
			"jwt.secret cannot be empty in production",
		)
	}

	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf(
			"jwt.expire_hours must be greater than 0",
		)
	}

	return nil
}
