package network

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// mockConn implements net.Conn interface for testing
type mockConn struct {
	readData  chan []byte
	writeData chan []byte
	closed    bool
}

func newMockConn() *mockConn {
	return &mockConn{
		readData:  make(chan []byte, 1000), // Increased buffer
		writeData: make(chan []byte, 1000), // Increased buffer for baud rate testing
	}
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	data := <-c.readData
	copy(b, data)
	return len(data), nil
}

func (c *mockConn) Write(b []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	c.writeData <- b
	return len(b), nil
}

func (c *mockConn) Close() error {
	c.closed = true
	return nil
}

func (c *mockConn) LocalAddr() net.Addr              { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0} }
func (c *mockConn) RemoteAddr() net.Addr             { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0} }
func (c *mockConn) SetDeadline(time.Time) error      { return nil }
func (c *mockConn) SetReadDeadline(time.Time) error  { return nil }
func (c *mockConn) SetWriteDeadline(time.Time) error { return nil }

// mockListener implements net.Listener interface for testing
type mockListener struct {
	connections chan net.Conn
	closed      bool
}

func newMockListener() *mockListener {
	return &mockListener{
		connections: make(chan net.Conn, 10),
	}
}

func (l *mockListener) Accept() (net.Conn, error) {
	if l.closed {
		return nil, io.EOF
	}
	conn := <-l.connections
	return conn, nil
}

func (l *mockListener) Close() error {
	l.closed = true
	return nil
}

func (l *mockListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
}

func TestNewTCPServer(t *testing.T) {
	cfg := Config{
		Host:           "localhost",
		Port:           10110,
		UpdateInterval: time.Second,
		Logger:         zerolog.Logger{},
		BaudRate:       4800,
	}

	server := NewTCPServer(cfg)

	if server == nil {
		t.Fatal("NewTCPServer returned nil")
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
}

// func TestTCPServerStartStop(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           10110,
// 		UpdateInterval: 100 * time.Millisecond,
// 		Logger:         zerolog.Logger{},
// 		BaudRate:       4800,
// 		SentenceOptions: SentenceOptions{
// 			EnablePosition: true, // Enable at least one type of sentence
// 		},
// 	}

// 	server := NewTCPServer(cfg)
// 	listener := newMockListener()
// 	server.listener = listener

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Start server in goroutine
// 	errCh := make(chan error, 1)
// 	go func() {
// 		errCh <- server.Start(ctx)
// 	}()

// 	// Add a mock client
// 	conn := newMockConn()
// 	listener.connections <- conn

// 	// Read data with timeout
// 	dataReceived := make(chan struct{})
// 	go func() {
// 		defer close(dataReceived)
// 		select {
// 		case data := <-conn.writeData:
// 			if strings.HasPrefix(string(data), "$") {
// 				return // Success
// 			}
// 		case <-time.After(1 * time.Second):
// 			t.Error("Timeout waiting for NMEA data")
// 		}
// 	}()

// 	<-dataReceived

// 	// Test graceful shutdown
// 	cancel()                           // Signal shutdown through context
// 	time.Sleep(100 * time.Millisecond) // Give time for cleanup

// 	if !listener.closed {
// 		t.Error("Listener not closed after Stop")
// 	}

// 	if !conn.closed {
// 		t.Error("Client connection not closed after Stop")
// 	}

// 	select {
// 	case err := <-errCh:
// 		if err != nil {
// 			t.Errorf("Server returned error: %v", err)
// 		}
// 	case <-time.After(500 * time.Millisecond):
// 		t.Error("Server did not shut down in time")
// 	}
// }

// func TestTCPServerBaudRate(t *testing.T) {
// 	cfg := Config{
// 		Host:           "localhost",
// 		Port:           10110,
// 		UpdateInterval: 100 * time.Millisecond,
// 		Logger:         zerolog.Logger{},
// 		BaudRate:       4800,
// 	}

// 	server := NewTCPServer(cfg)
// 	listener := newMockListener()
// 	server.listener = listener

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	go server.Start(ctx)

// 	conn := newMockConn()
// 	listener.connections <- conn

// 	// Wait for some data
// 	time.Sleep(200 * time.Millisecond)

// 	// Check that data is being sent at approximately the correct baud rate
// 	startTime := time.Now()
// 	bytesReceived := 0
// 	timeout := time.After(1 * time.Second)

// 	for {
// 		select {
// 		case data := <-conn.writeData:
// 			bytesReceived += len(data)
// 		case <-timeout:
// 			elapsed := time.Since(startTime)
// 			actualBytesPerSecond := float64(bytesReceived) / elapsed.Seconds()
// 			expectedBytesPerSecond := float64(cfg.BaudRate) / 10 // 10 bits per byte

// 			// Allow 20% margin of error due to testing environment variations
// 			margin := 0.2
// 			if actualBytesPerSecond < expectedBytesPerSecond*(1-margin) ||
// 				actualBytesPerSecond > expectedBytesPerSecond*(1+margin) {
// 				t.Errorf("Baud rate not within expected range. Got: %.2f bytes/s, Expected: %.2f bytes/s",
// 					actualBytesPerSecond, expectedBytesPerSecond)
// 			}
// 			return
// 		}
// 	}
// }
