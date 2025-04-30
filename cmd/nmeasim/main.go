// Package main provides the NMEA simulator command line interface
// supporting both NMEA 0183 and NMEA 2000 protocols
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/captv89/nmea-simulator/pkg/network"
	"github.com/captv89/nmea-simulator/pkg/nmea2000"
	"github.com/rs/zerolog"
)

func main() {
	// Command line flags
	protocol := flag.String("protocol", "both", "Protocol to use: nmea0183, nmea2000, or both")

	// NMEA 0183 flags
	nmea0183WSPort := flag.Int("nmea0183-ws-port", 8080, "WebSocket server port for NMEA 0183")
	nmea0183TCPPort := flag.Int("nmea0183-tcp-port", 10110, "TCP server port for NMEA 0183")
	baudRate := flag.Int("baud", 4800, "Baud rate for NMEA 0183 TCP output (4800, 9600, 19200, 38400)")

	// NMEA 2000 flags
	nmea2000WSPort := flag.Int("nmea2000-ws-port", 8081, "WebSocket server port for NMEA 2000")
	nmea2000TCPPort := flag.Int("nmea2000-tcp-port", 10200, "TCP port for NMEA 2000")

	// Common flags
	host := flag.String("host", "0.0.0.0", "Host to bind servers to")
	interval := flag.Duration("interval", time.Second, "Data update interval")
	flag.Parse()

	// Validate baud rate
	validBaudRates := map[int]bool{4800: true, 9600: true, 19200: true, 38400: true}
	if !validBaudRates[*baudRate] {
		*baudRate = 4800 // Default to 4800 if invalid
	}

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	var nmea0183Servers []network.Server
	var nmea2000Sim *nmea2000.Simulator

	// Start NMEA 0183 servers if needed
	if *protocol == "both" || *protocol == "nmea0183" {
		// Create NMEA 0183 server configuration
		cfg := network.Config{
			Host:           *host,
			UpdateInterval: *interval,
			Logger:         logger,
			BaudRate:       *baudRate,
			Protocol:       "nmea0183",
			SentenceOptions: network.SentenceOptions{
				EnablePosition:    true,
				EnableNavigation:  true,
				EnableEnvironment: true,
			},
		}

		// Create WebSocket server
		wsCfg := cfg
		wsCfg.Port = *nmea0183WSPort
		wsServer := network.NewWebSocketServer(wsCfg)
		nmea0183Servers = append(nmea0183Servers, wsServer)

		// Create TCP server
		tcpCfg := cfg
		tcpCfg.Port = *nmea0183TCPPort
		tcpServer := network.NewTCPServer(tcpCfg)
		nmea0183Servers = append(nmea0183Servers, tcpServer)

		// Start NMEA 0183 servers
		for _, server := range nmea0183Servers {
			srv := server // Create new variable for goroutine
			go func() {
				if err := srv.Start(ctx); err != nil {
					logger.Error().Err(err).Msg("NMEA 0183 server failed")
				}
			}()
		}

		logger.Info().
			Int("ws_port", *nmea0183WSPort).
			Int("tcp_port", *nmea0183TCPPort).
			Msg("NMEA 0183 simulator started")
	}

	// Start NMEA 2000 servers if needed
	if *protocol == "both" || *protocol == "nmea2000" {
		// Create NMEA 2000 TCP server
		tcpCfg := network.Config{
			Host:           *host,
			Port:           *nmea2000TCPPort,
			UpdateInterval: *interval,
			Logger:         logger,
			Protocol:       "nmea2000",
		}
		tcpServer := network.NewTCP2000Server(tcpCfg)

		// Create NMEA 2000 WebSocket server
		wsCfg := network.Config{
			Host:           *host,
			Port:           *nmea2000WSPort,
			UpdateInterval: *interval,
			Logger:         logger,
			Protocol:       "nmea2000",
		}
		wsServer := network.NewWebSocket2000Server(wsCfg)

		// Create and start NMEA 2000 simulator
		nmea2000Sim = nmea2000.New(nmea2000.Config{
			Transport:    tcpServer,
			WebSocket:    wsServer,
			UpdatePeriod: *interval,
		})

		// Start WebSocket server
		go func() {
			if err := wsServer.Start(ctx); err != nil {
				logger.Error().Err(err).Msg("NMEA 2000 websocket server failed")
			}
		}()

		// Start simulator (which starts TCP server)
		if err := nmea2000Sim.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("NMEA 2000 simulator failed to start")
			os.Exit(1)
		}

		logger.Info().
			Int("ws_port", *nmea2000WSPort).
			Int("tcp_port", *nmea2000TCPPort).
			Msg("NMEA 2000 simulator started")
	}

	if *protocol != "both" && *protocol != "nmea0183" && *protocol != "nmea2000" {
		logger.Error().Str("protocol", *protocol).Msg("invalid protocol specified")
		os.Exit(1)
	}

	// Wait for interrupt signal
	<-sigChan
	logger.Info().Msg("shutting down simulators...")

	// Stop NMEA 0183 servers
	for _, server := range nmea0183Servers {
		server.Stop()
	}

	// Stop NMEA 2000 simulator
	if nmea2000Sim != nil {
		nmea2000Sim.Stop()
	}
}
