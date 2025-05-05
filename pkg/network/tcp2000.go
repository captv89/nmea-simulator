package network

import (
	"context"
	"fmt"
	"net"

	"github.com/captv89/nmea-simulator/pkg/nmea2000/pgn"
)

// TCP2000Server implements NMEA 2000 message streaming over TCP
type TCP2000Server struct {
	*BaseServer
	listener net.Listener
	clients  map[net.Conn]bool
}

// NewTCP2000Server creates a new TCP server instance for NMEA 2000
func NewTCP2000Server(cfg Config) *TCP2000Server {
	return &TCP2000Server{
		BaseServer: NewBaseServer(cfg),
		clients:    make(map[net.Conn]bool),
	}
}

// Start begins the TCP server
func (s *TCP2000Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}
	s.listener = listener

	go s.acceptLoop(ctx)

	s.Config.Logger.Info().
		Str("addr", addr).
		Msg("NMEA 2000 TCP server started")

	return nil
}

// Stop terminates the TCP server
func (s *TCP2000Server) Stop() error {
	if s.listener != nil {
		s.listener.Close()
	}
	close(s.Done)

	s.Mu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[net.Conn]bool)
	s.Mu.Unlock()

	return nil
}

// SendPGN sends a NMEA 2000 message to all connected clients
func (s *TCP2000Server) SendPGN(msg pgn.Message) error {
	frame := formatPGNMessage(msg)
	var failedClients []net.Conn

	// Read lock for iterating
	s.Mu.RLock()
	for client := range s.clients {
		_, err := client.Write([]byte(frame))
		if err != nil {
			failedClients = append(failedClients, client)
		}
	}
	s.Mu.RUnlock()

	// Write lock for removing failed clients
	if len(failedClients) > 0 {
		s.Mu.Lock()
		for _, client := range failedClients {
			delete(s.clients, client)
			client.Close()
		}
		s.Mu.Unlock()
	}

	return nil
}

func (s *TCP2000Server) acceptLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				s.Config.Logger.Error().Err(err).Msg("Failed to accept connection")
				continue
			}
			s.Mu.Lock()
			s.clients[conn] = true
			s.Mu.Unlock()
			s.Config.Logger.Info().
				Str("remote", conn.RemoteAddr().String()).
				Msg("New NMEA 2000 client connected")
		}
	}
}

// formatPGNMessage formats a NMEA 2000 message for TCP transport
// Format: $PNMEA2K,PGN,Length,Data*Checksum
func formatPGNMessage(msg pgn.Message) string {
	var checksum byte
	data := fmt.Sprintf("$PNMEA2K,%d,%d,", msg.PGN, len(msg.Data))

	// Calculate checksum (XOR of all bytes after $ and before *)
	for i := 1; i < len(data); i++ {
		checksum ^= data[i]
	}
	for _, b := range msg.Data {
		checksum ^= b
	}

	return fmt.Sprintf("%s%X*%02X\r\n", data, msg.Data, checksum)
}
