// Package network provides network streaming capabilities for NMEA sentences
package network

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Server represents a network server interface that can stream NMEA sentences
type Server interface {
	Start(ctx context.Context) error
	Stop() error
}

// Config holds server configuration
type Config struct {
	Host            string
	Port            int
	UpdateInterval  time.Duration
	Logger          zerolog.Logger
	SentenceOptions SentenceOptions
	BaudRate        int // Added baud rate configuration
}

// SentenceOptions configures which NMEA sentences to generate
type SentenceOptions struct {
	EnablePosition    bool // GGA, GLL
	EnableNavigation  bool // RMC, HDT, VTG, XTE
	EnableEnvironment bool // DBT, MTW, MWV, VHW, DPT
}

// BaseServer provides common functionality for TCP and WebSocket servers
type BaseServer struct {
	Config Config
	Mu     sync.RWMutex
	Done   chan struct{}
}

// NewBaseServer creates a new base server with the given configuration
func NewBaseServer(cfg Config) *BaseServer {
	return &BaseServer{
		Config: cfg,
		Done:   make(chan struct{}),
	}
}
