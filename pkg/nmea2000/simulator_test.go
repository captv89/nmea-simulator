package nmea2000

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea2000/pgn"
)

// mockNMEA2000Server implements network.NMEA2000Server for testing
type mockNMEA2000Server struct {
	mu      sync.Mutex
	started bool
	stopped bool
	msgs    []pgn.Message
}

func (m *mockNMEA2000Server) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.started = true
	return nil
}

func (m *mockNMEA2000Server) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = true
	return nil
}

func (m *mockNMEA2000Server) SendPGN(msg pgn.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, msg)
	return nil
}

func (m *mockNMEA2000Server) getMessages() []pgn.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.msgs
}

func TestSimulator(t *testing.T) {
	transport := &mockNMEA2000Server{}
	webSocket := &mockNMEA2000Server{}

	sim := New(Config{
		Transport:    transport,
		WebSocket:    webSocket,
		UpdatePeriod: 100 * time.Millisecond,
	})

	// Start simulator
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := sim.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start simulator: %v", err)
	}

	// Wait for a few update cycles
	time.Sleep(220 * time.Millisecond)

	// Stop simulator
	if err := sim.Stop(); err != nil {
		t.Fatalf("Failed to stop simulator: %v", err)
	}

	// Verify transport server state
	if !transport.started {
		t.Error("Transport server was not started")
	}
	if !transport.stopped {
		t.Error("Transport server was not stopped")
	}

	// Verify websocket server state
	if !webSocket.started {
		t.Error("WebSocket server was not started")
	}
	if !webSocket.stopped {
		t.Error("WebSocket server was not stopped")
	}

	// Check messages were sent to both servers
	transportMsgs := transport.getMessages()
	webSocketMsgs := webSocket.getMessages()

	if len(transportMsgs) == 0 {
		t.Error("No messages sent to transport server")
	}
	if len(webSocketMsgs) == 0 {
		t.Error("No messages sent to WebSocket server")
	}

	// Verify message types
	pgns := make(map[uint32]bool)
	for _, msg := range transportMsgs {
		pgns[msg.PGN] = true
	}

	expectedPGNs := []uint32{
		127250, // Vessel Heading
		128259, // Speed
		128267, // Water Depth
		129025, // Position
		130306, // Wind Data
	}

	for _, expected := range expectedPGNs {
		if !pgns[expected] {
			t.Errorf("Expected PGN %d not found in messages", expected)
		}
	}
}

func TestSimulatorWithoutWebSocket(t *testing.T) {
	transport := &mockNMEA2000Server{}

	sim := New(Config{
		Transport:    transport,
		WebSocket:    nil, // No WebSocket server
		UpdatePeriod: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := sim.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start simulator: %v", err)
	}

	time.Sleep(220 * time.Millisecond)

	if err := sim.Stop(); err != nil {
		t.Fatalf("Failed to stop simulator: %v", err)
	}

	if !transport.started {
		t.Error("Transport server was not started")
	}
	if !transport.stopped {
		t.Error("Transport server was not stopped")
	}

	msgs := transport.getMessages()
	if len(msgs) == 0 {
		t.Error("No messages sent to transport server")
	}
}

func TestSimulatorContextCancellation(t *testing.T) {
	transport := &mockNMEA2000Server{}
	webSocket := &mockNMEA2000Server{}

	sim := New(Config{
		Transport:    transport,
		WebSocket:    webSocket,
		UpdatePeriod: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())

	err := sim.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start simulator: %v", err)
	}

	// Let it run briefly
	time.Sleep(150 * time.Millisecond)

	// Cancel context
	cancel()

	// Give it time to shut down
	time.Sleep(150 * time.Millisecond)

	msgs := transport.getMessages()
	lastCount := len(msgs)

	// Verify no more messages are being sent
	time.Sleep(150 * time.Millisecond)
	msgs = transport.getMessages()
	if len(msgs) > lastCount {
		t.Error("Simulator continued sending messages after context cancellation")
	}
}