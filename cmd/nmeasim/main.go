package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/captv89/nmea-simulator/pkg/network"
	"github.com/rs/zerolog"
)

func main() {
	// Command line flags
	wsPort := flag.Int("ws-port", 8080, "WebSocket server port")
	tcpPort := flag.Int("tcp-port", 10110, "TCP server port (default NMEA port)")
	host := flag.String("host", "0.0.0.0", "Host to bind servers to")
	interval := flag.Duration("interval", time.Second, "NMEA sentence update interval")
	baudRate := flag.Int("baud", 4800, "Baud rate for TCP output (4800, 9600, 19200, 38400)")
	flag.Parse()

	// Validate baud rate
	validBaudRates := map[int]bool{4800: true, 9600: true, 19200: true, 38400: true}
	if !validBaudRates[*baudRate] {
		*baudRate = 4800 // Default to 4800 if invalid
	}

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create server configuration
	cfg := network.Config{
		Host:           *host,
		UpdateInterval: *interval,
		Logger:         logger,
		BaudRate:       *baudRate,
		SentenceOptions: network.SentenceOptions{
			EnablePosition:    true,
			EnableNavigation:  true,
			EnableEnvironment: true,
		},
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WebSocket server
	wsCfg := cfg
	wsCfg.Port = *wsPort
	wsServer := network.NewWebSocketServer(wsCfg)

	// Create TCP server
	tcpCfg := cfg
	tcpCfg.Port = *tcpPort
	tcpServer := network.NewTCPServer(tcpCfg)

	// Start servers
	go func() {
		if err := wsServer.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("websocket server failed")
		}
	}()

	go func() {
		if err := tcpServer.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("tcp server failed")
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	logger.Info().Msg("shutting down servers...")

	cancel()
	wsServer.Stop()
	tcpServer.Stop()
}
