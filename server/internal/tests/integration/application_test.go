package integration

import "testing"

func TestApplication_Start(t *testing.T) {
	app := NewTestApplication(t)

	if app == nil {
		t.Fatal("expected test application")
	}

	if app.Server.HTTPHandler() == nil {
		t.Fatal("expected router")
	}

	if app.Server.DB() == nil {
		t.Fatal("expected database")
	}
}
