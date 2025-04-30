package network

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

func TestNewWebSocket2000Server(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8081,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
		Protocol:       "nmea2000",
	}

	server := NewWebSocket2000Server(cfg)

	if server == nil {
		t.Fatal("NewWebSocket2000Server returned nil")
	}

	// Check fields
	if server.Config.Host != cfg.Host {
		t.Error("Host not set correctly")
	}
	if server.Config.Port != cfg.Port {
		t.Error("Port not set correctly")
	}
	if server.Config.UpdateInterval != cfg.UpdateInterval {
		t.Error("UpdateInterval not set correctly")
	}
	if server.Config.Protocol != cfg.Protocol {
		t.Error("Protocol not set correctly")
	}

	if server.clients == nil {
		t.Error("Clients map not initialized")
	}
}

// func TestWebSocket2000ServerConnection(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           8081,
// 		UpdateInterval: 100 * time.Millisecond,
// 		Logger:         zerolog.Logger{},
// 		Protocol:       "nmea2000",
// 	}

// 	server := NewWebSocket2000Server(cfg)

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
// 	defer ws.Close()

// 	// Test sending PGN message
// 	testPGN := pgn.Message{
// 		PGN:  128259, // Speed
// 		Data: []byte{0x01, 0x02, 0x03},
// 	}

// 	if err := server.SendPGN(testPGN); err != nil {
// 		t.Fatalf("Failed to send PGN: %v", err)
// 	}

// 	// Check that message was received
// 	timeout := time.After(time.Second)
// 	messageReceived := false
// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)
// 		_, message, err := ws.ReadMessage()
// 		if err != nil {
// 			t.Errorf("Failed to read WebSocket message: %v", err)
// 			return
// 		}

// 		// Verify NMEA 2000 message format
// 		msg := string(message)
// 		if !strings.HasPrefix(msg, "$PNMEA2K") {
// 			t.Error("Received message is not a valid NMEA 2000 message")
// 			return
// 		}
// 		if !strings.Contains(msg, "128259") {
// 			t.Error("PGN not found in message")
// 			return
// 		}
// 		messageReceived = true
// 	}()

// 	select {
// 	case <-timeout:
// 		t.Error("Timeout waiting for WebSocket message")
// 	case <-done:
// 		if !messageReceived {
// 			t.Error("No valid NMEA 2000 message received")
// 		}
// 	}
// }

func TestWebSocket2000ServerClientManagement(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           8081,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
		Protocol:       "nmea2000",
	}

	server := NewWebSocket2000Server(cfg)

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

// func TestWebSocket2000ServerGracefulShutdown(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           8081,
// 		UpdateInterval: time.Second,
// 		Logger:         zerolog.Logger{},
// 		Protocol:       "nmea2000",
// 	}

// 	server := NewWebSocket2000Server(cfg)

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
