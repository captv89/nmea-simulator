package network

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestBaseServer_Config(t *testing.T) {
	cfg := Config{
		Host: 		 "localhost",
		Port: 		 8080,
		Logger:  zerolog.New(os.Stdout),
	}
	server := NewBaseServer(cfg)

	// Compare fields individually since Logger can't be directly compared
	if server.Config.Host != cfg.Host {
		t.Errorf("expected host %s, got %s", cfg.Host, server.Config.Host)
	}
	if server.Config.Port != cfg.Port {
		t.Errorf("expected port %d, got %d", cfg.Port, server.Config.Port)
	}
}
