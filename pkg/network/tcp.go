package network

import (
	"context"
	"fmt"
	"net"
	"strings"
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

	go s.acceptLoop(ctx)
	go s.broadcastLoop(ctx)

	<-ctx.Done()
	return s.Stop()
}

func (s *TCPServer) acceptLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Done:
			return
		default:
			if conn, err := s.listener.Accept(); err == nil {
				s.Mu.Lock()
				s.clients[conn] = true
				s.Mu.Unlock()

				s.Config.Logger.Info().
					Str("remote", conn.RemoteAddr().String()).
					Msg("new TCP client connected")

				// Monitor connection for closure
				go func(conn net.Conn) {
					select {
					case <-s.Done:
						conn.Close()
					case <-ctx.Done():
						conn.Close()
					}
				}(conn)
			} else if !strings.Contains(err.Error(), "use of closed network connection") {
				s.Config.Logger.Error().Err(err).Msg("accept error")
			}
		}
	}
}

func (s *TCPServer) broadcastLoop(ctx context.Context) {
	ticker := time.NewTicker(s.Config.UpdateInterval)
	defer ticker.Stop()

	// Calculate bytes per interval based on baud rate
	bytesPerInterval := int(float64(s.Config.BaudRate) * s.Config.UpdateInterval.Seconds() / 8)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Done:
			return
		case <-ticker.C:
			sentences := s.generateSentences()
			s.broadcast(sentences, bytesPerInterval)
		}
	}
}

func (s *TCPServer) broadcast(sentences []string, bytesPerInterval int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var totalBytes int
	for conn := range s.clients {
		for _, sentence := range sentences {
			if totalBytes >= bytesPerInterval {
				return // Respect baud rate limit
			}

			data := []byte(sentence + "\r\n")
			totalBytes += len(data)

			_, err := conn.Write(data)
			if err != nil {
				s.Config.Logger.Error().
					Err(err).
					Str("remote", conn.RemoteAddr().String()).
					Msg("failed to send message")
				conn.Close()
				delete(s.clients, conn)
				break
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

// Stop closes all client connections and stops the server
func (s *TCPServer) Stop() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	select {
	case <-s.Done:
		return nil
	default:
		close(s.Done)
	}

	// Close all client connections
	for conn := range s.clients {
		conn.Close()
		delete(s.clients, conn)
	}

	// Close listener
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}

	return nil
}
