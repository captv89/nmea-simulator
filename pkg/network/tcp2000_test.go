package network

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea2000/pgn"
	"github.com/rs/zerolog"
)

// mockTCP2000Conn is a mock implementation of net.Conn for testing
type mockTCP2000Conn struct {
	readData  chan []byte
	writeData chan []byte
	closed    bool
}

func newMockTCP2000Conn() *mockTCP2000Conn {
	return &mockTCP2000Conn{
		readData:  make(chan []byte, 100),
		writeData: make(chan []byte, 100),
	}
}

func (c *mockTCP2000Conn) Read(b []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	data := <-c.readData
	copy(b, data)
	return len(data), nil
}

func (c *mockTCP2000Conn) Write(b []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	c.writeData <- b
	return len(b), nil
}

func (c *mockTCP2000Conn) Close() error {
	c.closed = true
	return nil
}

func (c *mockTCP2000Conn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func (c *mockTCP2000Conn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}
func (c *mockTCP2000Conn) SetDeadline(t time.Time) error      { return nil }
func (c *mockTCP2000Conn) SetReadDeadline(t time.Time) error  { return nil }
func (c *mockTCP2000Conn) SetWriteDeadline(t time.Time) error { return nil }

// mockTCP2000Listener is a mock implementation of net.Listener for testing
type mockTCP2000Listener struct {
	connections chan net.Conn
	closed      bool
}

func newMockTCP2000Listener() *mockTCP2000Listener {
	return &mockTCP2000Listener{
		connections: make(chan net.Conn, 10),
	}
}

func (l *mockTCP2000Listener) Accept() (net.Conn, error) {
	if l.closed {
		return nil, io.EOF
	}
	conn := <-l.connections
	return conn, nil
}

func (l *mockTCP2000Listener) Close() error {
	l.closed = true
	return nil
}

func (l *mockTCP2000Listener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func TestNewTCP2000Server(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           10200,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
		Protocol:       "nmea2000",
	}

	server := NewTCP2000Server(cfg)

	if server == nil {
		t.Fatal("NewTCP2000Server returned nil")
	}

	if server.Config.Host != cfg.Host {
		t.Error("Host not set correctly")
	}
	if server.Config.Port != cfg.Port {
		t.Error("Port not set correctly")
	}
	if server.Config.Protocol != "nmea2000" {
		t.Error("Protocol not set correctly")
	}
	if server.clients == nil {
		t.Error("Clients map not initialized")
	}
}

// func TestTCP2000ServerConnection(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           10200,
// 		UpdateInterval: 100 * time.Millisecond,
// 		Logger:         zerolog.Logger{},
// 		Protocol:       "nmea2000",
// 	}

// 	server := NewTCP2000Server(cfg)
// 	listener := newMockTCP2000Listener()
// 	server.listener = listener

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Start server
// 	go func() {
// 		if err := server.Start(ctx); err != nil {
// 			t.Errorf("Server failed to start: %v", err)
// 		}
// 	}()

// 	// Create and add a mock client
// 	conn := newMockTCP2000Conn()
// 	listener.connections <- conn

// 	// Test sending PGN message
// 	testPGN := pgn.Message{
// 		PGN:  128259, // Speed
// 		Data: []byte{0x01, 0x02, 0x03},
// 	}

// 	if err := server.SendPGN(testPGN); err != nil {
// 		t.Fatalf("Failed to send PGN: %v", err)
// 	}

// 	// Read data with timeout
// 	select {
// 	case data := <-conn.writeData:
// 		msg := string(data)
// 		if !strings.HasPrefix(msg, "$PNMEA2K") {
// 			t.Error("Message does not have correct prefix")
// 		}
// 		if !strings.Contains(msg, "128259") {
// 			t.Error("PGN not found in message")
// 		}
// 	case <-time.After(time.Second):
// 		t.Error("Timeout waiting for PGN message")
// 	}
// }

// func TestTCP2000ServerGracefulShutdown(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           10200,
// 		UpdateInterval: time.Second,
// 		Logger:         zerolog.Logger{},
// 		Protocol:       "nmea2000",
// 	}

// 	server := NewTCP2000Server(cfg)
// 	listener := newMockTCP2000Listener()
// 	server.listener = listener

// 	// Add some mock clients
// 	conn1 := newMockTCP2000Conn()
// 	conn2 := newMockTCP2000Conn()

// 	server.Mu.Lock()
// 	server.clients[conn1] = true
// 	server.clients[conn2] = true
// 	server.Mu.Unlock()

// 	// Test graceful shutdown
// 	if err := server.Stop(); err != nil {
// 		t.Errorf("Error during shutdown: %v", err)
// 	}

// 	// Verify listener is closed
// 	if !listener.closed {
// 		t.Error("Listener not closed after shutdown")
// 	}

// 	// Verify clients are closed
// 	if !conn1.closed || !conn2.closed {
// 		t.Error("Not all clients were closed during shutdown")
// 	}

// 	// Verify clients map is empty
// 	server.Mu.Lock()
// 	clientCount := len(server.clients)
// 	server.Mu.Unlock()

// 	if clientCount != 0 {
// 		t.Errorf("Expected 0 clients after shutdown, got %d", clientCount)
// 	}
// }

func TestFormatPGNMessage(t *testing.T) {
	testCases := []struct {
		name     string
		msg      pgn.Message
		wantPGN  string
		wantData string
	}{
		{
			name: "Speed message",
			msg: pgn.Message{
				PGN:  128259,
				Data: []byte{0x01, 0x02, 0x03},
			},
			wantPGN:  "128259",
			wantData: "010203",
		},
		{
			name: "Depth message",
			msg: pgn.Message{
				PGN:  128267,
				Data: []byte{0xFF},
			},
			wantPGN:  "128267",
			wantData: "FF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatPGNMessage(tc.msg)

			if !strings.Contains(result, "$PNMEA2K") {
				t.Error("Missing message prefix")
			}
			if !strings.Contains(result, tc.wantPGN) {
				t.Error("PGN not found in message")
			}
			if !strings.Contains(result, tc.wantData) {
				t.Error("Data not properly formatted")
			}
			if !strings.HasSuffix(result, "\r\n") {
				t.Error("Missing message terminator")
			}
		})
	}
}
