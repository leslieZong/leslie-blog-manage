package integration

import (
	"leslie-blog-server/internal/bootstrap"
	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/server"
	"testing"
)

type TestApplication struct {
	*bootstrap.Server
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

	if cfg.App.Env != "test" {
		t.Fatalf(
			"refuse to run integration test with env=%q",
			cfg.App.Env,
		)
	}

	return cfg
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
	}
}
