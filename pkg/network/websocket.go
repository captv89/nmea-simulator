package network

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/environment"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/navigation"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/position"
	"github.com/gorilla/websocket"
)

//go:embed web/*
var webContent embed.FS

// WebSocketServer implements NMEA sentence streaming over WebSocket
type WebSocketServer struct {
	*BaseServer
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	clientMu sync.Mutex
}

// NewWebSocketServer creates a new WebSocket server instance
func NewWebSocketServer(cfg Config) *WebSocketServer {
	return &WebSocketServer{
		BaseServer: NewBaseServer(cfg),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins for testing
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// loggingMiddleware wraps an http.Handler and logs requests
func (s *WebSocketServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.Config.Logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote", r.RemoteAddr).
			Msg("incoming request")

		next.ServeHTTP(w, r)

		s.Config.Logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote", r.RemoteAddr).
			Dur("duration", time.Since(start)).
			Msg("request completed")
	})
}

// Start begins the WebSocket server
func (s *WebSocketServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	// Setup routes
	mux := http.NewServeMux()

	// Serve static files from the embedded filesystem
	fsys, err := fs.Sub(webContent, "web")
	if err != nil {
		return fmt.Errorf("failed to setup static file serving: %w", err)
	}

	// Create a stripped file system to serve index.html from root
	fsRoot := http.FileServer(http.FS(fsys))

	// Handle root path specifically
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Directly serve index.html for root path
			data, err := fs.ReadFile(fsys, "index.html")
			if err != nil {
				s.Config.Logger.Error().Err(err).Msg("failed to read index.html")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
			return
		}

		// For all other paths, use the file server
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		fsRoot.ServeHTTP(w, r)
	})

	// Handle WebSocket path
	mux.HandleFunc("/ws", s.handleWebSocket)

	// Create server with logging middleware
	handler := s.loggingMiddleware(mux)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	s.Config.Logger.Info().Str("addr", addr).Msg("starting websocket server")

	// Handle server shutdown
	go func() {
		<-ctx.Done()
		s.Config.Logger.Info().Msg("shutting down websocket server")
		server.Close()
	}()

	// Start periodic data transmission
	go s.broadcastLoop(ctx)

	return server.ListenAndServe()
}

// handleWebSocket handles WebSocket connections
func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Config.Logger.Error().Err(err).Msg("websocket upgrade failed")
		return
	}

	// Register new client
	s.clientMu.Lock()
	s.clients[conn] = true
	s.clientMu.Unlock()

	s.Config.Logger.Info().Str("remote", conn.RemoteAddr().String()).Msg("new client connected")

	// Create a done channel for this connection
	done := make(chan struct{})

	// Handle client messages and connection status
	go func() {
		defer func() {
			s.clientMu.Lock()
			delete(s.clients, conn)
			s.clientMu.Unlock()
			conn.Close()
			close(done)
			s.Config.Logger.Info().Str("remote", conn.RemoteAddr().String()).Msg("client disconnected")
		}()

		for {
			// Read messages from client (if any)
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					s.Config.Logger.Error().Err(err).Str("remote", conn.RemoteAddr().String()).Msg("websocket error")
				}
				return
			}
		}
	}()

	// Wait for client disconnection or context cancellation
	select {
	case <-r.Context().Done():
	case <-s.Done:
	case <-done:
	}
}

// Stop closes all client connections and stops the server
func (s *WebSocketServer) Stop() error {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	for client := range s.clients {
		client.Close()
		delete(s.clients, client)
	}

	close(s.Done)
	return nil
}

func (s *WebSocketServer) broadcastLoop(ctx context.Context) {
	ticker := s.Config.UpdateInterval
	t := time.NewTicker(ticker)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Done:
			return
		case <-t.C:
			sentences := s.generateSentences()
			s.broadcast(sentences)
		}
	}
}

func (s *WebSocketServer) generateSentences() []string {
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

func (s *WebSocketServer) broadcast(sentences []string) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	for client := range s.clients {
		for _, sentence := range sentences {
			err := client.WriteMessage(websocket.TextMessage, []byte(sentence))
			if err != nil {
				s.Config.Logger.Error().Err(err).
					Str("remote", client.RemoteAddr().String()).
					Msg("failed to send message")
				client.Close()
				delete(s.clients, client)
				break
			}
		}
	}
}
