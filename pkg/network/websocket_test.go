package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

func TestNewWebSocketServer(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8080,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
		BaudRate:       4800,
	}

	server := NewWebSocketServer(cfg)

	if server == nil {
		t.Fatal("NewWebSocketServer returned nil")
	}

	// Check individual fields instead of comparing entire struct
	if server.Config.Host != cfg.Host {
		t.Error("Host not set correctly")
	}
	if server.Config.Port != cfg.Port {
		t.Error("Port not set correctly")
	}
	if server.Config.UpdateInterval != cfg.UpdateInterval {
		t.Error("UpdateInterval not set correctly")
	}
	if server.Config.BaudRate != cfg.BaudRate {
		t.Error("BaudRate not set correctly")
	}

	if server.clients == nil {
		t.Error("Clients map not initialized")
	}

	if server.upgrader.CheckOrigin == nil {
		t.Error("WebSocket upgrader not configured properly")
	}
}

func TestWebSocketServerConnection(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8080,
		UpdateInterval: 100 * time.Millisecond,
		Logger:         zerolog.Logger{},
		SentenceOptions: SentenceOptions{
			EnablePosition:    true,
			EnableNavigation:  true,
			EnableEnvironment: true,
		},
	}

	server := NewWebSocketServer(cfg)

	// Create test HTTP server
	s := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer s.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not connect to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Start the broadcast loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.broadcastLoop(ctx)

	// Check that we receive NMEA sentences
	messageReceived := false
	timeout := time.After(2 * time.Second)
	done := make(chan struct{})

	go func() {
		defer close(done)
		_, message, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("Failed to read WebSocket message: %v", err)
			return
		}
		if !strings.HasPrefix(string(message), "$") {
			t.Error("Received message is not a valid NMEA sentence")
			return
		}
		messageReceived = true
	}()

	select {
	case <-timeout:
		t.Error("Timeout waiting for WebSocket message")
	case <-done:
		if !messageReceived {
			t.Error("No valid NMEA message received")
		}
	}
}

func TestWebSocketServerStaticFiles(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8080,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
	}

	server := NewWebSocketServer(cfg)

	// Create test HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/ws") {
			server.handleWebSocket(w, r)
			return
		}

		// For root and static files, simulate the same behavior as in Start()
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<!DOCTYPE html><html><body>Test Page</body></html>"))
			return
		}
	})

	s := httptest.NewServer(server.loggingMiddleware(handler))
	defer s.Close()

	// Test root path
	resp, err := http.Get(s.URL + "/")
	if err != nil {
		t.Fatalf("Failed to get root path: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}
}

// func TestWebSocketServerGracefulShutdown(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           8080,
// 		UpdateInterval: time.Second,
// 		Logger:         zerolog.Logger{},
// 	}

// 	server := NewWebSocketServer(cfg)

// 	// Create test HTTP server
// 	s := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
// 	defer s.Close()

// 	// Convert http URL to ws URL
// 	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

// 	// Connect WebSocket client
// 	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		t.Fatalf("Could not connect to WebSocket server: %v", err)
// 	}

// 	// Test graceful shutdown
// 	server.Stop()

// 	// Try to send a message, should fail
// 	err = ws.WriteMessage(websocket.TextMessage, []byte("test"))
// 	if err == nil {
// 		t.Error("Expected connection to be closed after shutdown")
// 	}
// }

func TestWebSocketServerClientManagement(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8080,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
	}

	server := NewWebSocketServer(cfg)

	// Create test HTTP server
	s := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer s.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect multiple clients
	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Could not connect client %d: %v", i, err)
		}
		clients[i] = ws
	}

	// Verify client count
	server.clientMu.Lock()
	clientCount := len(server.clients)
	server.clientMu.Unlock()

	if clientCount != 3 {
		t.Errorf("Expected 3 clients, got %d", clientCount)
	}

	// Disconnect one client
	clients[0].Close()
	time.Sleep(100 * time.Millisecond)

	// Verify client was removed
	server.clientMu.Lock()
	clientCount = len(server.clients)
	server.clientMu.Unlock()

	if clientCount != 2 {
		t.Errorf("Expected 2 clients after disconnect, got %d", clientCount)
	}

	// Cleanup remaining clients
	for i := 1; i < 3; i++ {
		clients[i].Close()
	}
}
