package integration

import (
	"leslie-blog-server/internal/bootstrap"
	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/server"
	"testing"
)

type TestApplication struct {
	*bootstrap.Server
	Close func() error
}

func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()

	cfg, err := config.Load("configs/config.test.yaml", ".env.test")
	if err != nil {
		t.Fatalf(
			"load test config failed: %v",
			err,
		)
	}

	validateTestEnvironment(t, cfg)

	return cfg
}
func validateTestEnvironment(
	t *testing.T,
	cfg *config.Config,
) {
	t.Helper()

	if cfg.App.Env != "test" {
		t.Fatalf(
			"refuse to run integration tests with env=%q",
			cfg.App.Env,
		)
	}

	if cfg.MySQL.Database != "leslie_blog_test" {
		t.Fatalf(
			"refuse to run integration tests with mysql database=%q",
			cfg.MySQL.Database,
		)
	}

	if cfg.Redis.DB != 1 {
		t.Fatalf(
			"refuse to run integration tests with redis db=%d",
			cfg.Redis.DB,
		)
	}

}

func NewTestApplication(t *testing.T) *TestApplication {
	t.Helper()
	cfg := loadTestConfig(t)
	srv, err := server.New(cfg)
	if err != nil {
		t.Fatalf(
			"create server failed: %v",
			err,
		)
	}

	return &TestApplication{
		Server: srv,
		Close:  srv.Close,
	}
}
