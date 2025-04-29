package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/environment"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/navigation"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/position"
)

// TCPServer implements NMEA sentence streaming over TCP
type TCPServer struct {
	*BaseServer
	listener net.Listener
	clients  map[net.Conn]bool
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(cfg Config) *TCPServer {
	return &TCPServer{
		BaseServer: NewBaseServer(cfg),
		clients:    make(map[net.Conn]bool),
	}
}

// Start begins the TCP server
func (s *TCPServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	s.listener = listener

	s.Config.Logger.Info().Str("addr", addr).Msg("starting TCP server")

	// Handle client connections
	go s.acceptConnections(ctx)

	// Handle server shutdown
	go func() {
		<-ctx.Done()
		s.Config.Logger.Info().Msg("shutting down TCP server")
		s.Stop()
	}()

	return nil
}

// Stop closes the listener and all client connections
func (s *TCPServer) Stop() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if s.listener != nil {
		s.listener.Close()
	}

	for client := range s.clients {
		client.Close()
		delete(s.clients, client)
	}

	close(s.Done)
	return nil
}

func (s *TCPServer) acceptConnections(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				s.Config.Logger.Error().Err(err).Msg("failed to accept connection")
				continue
			}
		}

		s.Mu.Lock()
		s.clients[conn] = true
		s.Mu.Unlock()

		s.Config.Logger.Info().Str("remote", conn.RemoteAddr().String()).Msg("new client connected")

		// Handle individual client in a goroutine
		go s.handleClient(ctx, conn)
	}
}

func (s *TCPServer) handleClient(ctx context.Context, conn net.Conn) {
	defer func() {
		s.Mu.Lock()
		delete(s.clients, conn)
		s.Mu.Unlock()
		conn.Close()
		s.Config.Logger.Info().Str("remote", conn.RemoteAddr().String()).Msg("client disconnected")
	}()

	ticker := time.NewTicker(s.Config.UpdateInterval)
	defer ticker.Stop()

	// Calculate bytes per second based on baud rate
	bytesPerSecond := s.Config.BaudRate / 10 // 8 data bits + 1 start bit + 1 stop bit = 10 bits per byte
	if bytesPerSecond == 0 {
		bytesPerSecond = 4800 / 10 // Default to 4800 baud if not set
	}

	// Create a rate limiter based on baud rate
	limiter := time.NewTicker(time.Second / time.Duration(bytesPerSecond))
	defer limiter.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Done:
			return
		case <-ticker.C:
			sentences := s.generateSentences()
			for _, sentence := range sentences {
				// Add CRLF as per NMEA spec
				sentence = sentence + "\r\n"
				data := []byte(sentence)

				// Send each byte at the configured baud rate
				for _, b := range data {
					<-limiter.C // Wait for the next tick before sending byte
					_, err := conn.Write([]byte{b})
					if err != nil {
						s.Config.Logger.Error().Err(err).
							Str("remote", conn.RemoteAddr().String()).
							Msg("failed to write to client")
						return
					}
				}
			}
		}
	}
}

func (s *TCPServer) generateSentences() []string {
	var sentences []string

	if s.Config.SentenceOptions.EnablePosition {
		sentences = append(sentences,
			position.GenerateGGA(),
			position.GenerateGLL(),
		)
	}

	if s.Config.SentenceOptions.EnableNavigation {
		sentences = append(sentences,
			navigation.GenerateRMC(),
			navigation.GenerateHDT(),
			navigation.GenerateVTG(),
			navigation.GenerateXTE(),
		)
	}

	if s.Config.SentenceOptions.EnableEnvironment {
		sentences = append(sentences,
			environment.GenerateDBT(),
			environment.GenerateMTW(),
			environment.GenerateMWV(),
			environment.GenerateVHW(),
			environment.GenerateDPT(),
		)
	}

	return sentences
}
