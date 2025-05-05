// Package network provides network streaming capabilities for NMEA sentences
package network

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea2000/pgn"
	"github.com/gorilla/websocket"
)

// WebSocket2000Server implements NMEA 2000 message streaming over WebSocket
type WebSocket2000Server struct {
	*BaseServer
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	clientMu sync.Mutex
}

// NewWebSocket2000Server creates a new WebSocket server instance for NMEA 2000
func NewWebSocket2000Server(cfg Config) *WebSocket2000Server {
	return &WebSocket2000Server{
		BaseServer: NewBaseServer(cfg),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start begins the WebSocket server
func (s *WebSocket2000Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	// Setup routes
	mux := http.NewServeMux()

	// Serve static files from the embedded filesystem
	fsys, err := fs.Sub(webContent, "web")
	if err != nil {
		return fmt.Errorf("failed to setup static file serving: %w", err)
	}

	// Create a file server for static files
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

	// Handle WebSocket path for NMEA 2000
	mux.HandleFunc("/nmea2000", s.handleWebSocket)

	// Create server with logging middleware
	handler := s.loggingMiddleware(mux)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	s.Config.Logger.Info().Str("addr", addr).Msg("starting NMEA 2000 websocket server")

	// Handle server shutdown
	go func() {
		<-ctx.Done()
		s.Config.Logger.Info().Msg("shutting down NMEA 2000 websocket server")
		server.Close()
	}()

	return server.ListenAndServe()
}

// Stop terminates the WebSocket server
func (s *WebSocket2000Server) Stop() error {
	close(s.Done)
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[*websocket.Conn]bool)
	return nil
}

// SendPGN sends a NMEA 2000 message to all connected WebSocket clients
func (s *WebSocket2000Server) SendPGN(msg pgn.Message) error {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	frame := formatPGNMessage(msg)

	for client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(frame))
		if err != nil {
			s.Config.Logger.Error().
				Err(err).
				Str("remote", client.RemoteAddr().String()).
				Msg("Failed to send message")

			client.Close()
			delete(s.clients, client)
		}
	}
	return nil
}

func (s *WebSocket2000Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Config.Logger.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}

	s.clientMu.Lock()
	s.clients[conn] = true
	s.clientMu.Unlock()

	s.Config.Logger.Info().
		Str("remote", conn.RemoteAddr().String()).
		Msg("New NMEA 2000 WebSocket client connected")

	// Create a done channel for this connection
	done := make(chan struct{})

	go func() {
		defer func() {
			s.clientMu.Lock()
			delete(s.clients, conn)
			s.clientMu.Unlock()
			conn.Close()
			close(done)
			s.Config.Logger.Info().
				Str("remote", conn.RemoteAddr().String()).
				Msg("NMEA 2000 client disconnected")
		}()

		for {
			// Read messages from client (if any)
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					s.Config.Logger.Error().Err(err).Msg("WebSocket read error")
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

// loggingMiddleware wraps an http.Handler and logs requests
func (s *WebSocket2000Server) loggingMiddleware(next http.Handler) http.Handler {
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
