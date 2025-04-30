// Package nmea2000 provides NMEA 2000 simulation capabilities
package nmea2000

import (
	"context"
	"time"

	"github.com/captv89/nmea-simulator/pkg/network"
	"github.com/captv89/nmea-simulator/pkg/nmea2000/pgn"
)

// Simulator represents a NMEA 2000 network simulator
type Simulator struct {
	transport    network.NMEA2000Server
	webSocket    network.NMEA2000Server
	updatePeriod time.Duration
	done         chan struct{}
}

// Config holds simulator configuration
type Config struct {
	Transport    network.NMEA2000Server
	WebSocket    network.NMEA2000Server
	UpdatePeriod time.Duration
}

// New creates a new NMEA 2000 simulator
func New(cfg Config) *Simulator {
	return &Simulator{
		transport:    cfg.Transport,
		webSocket:    cfg.WebSocket,
		updatePeriod: cfg.UpdatePeriod,
		done:         make(chan struct{}),
	}
}

// Start begins the simulation
func (s *Simulator) Start(ctx context.Context) error {
	if err := s.transport.Start(ctx); err != nil {
		return err
	}

	if s.webSocket != nil {
		if err := s.webSocket.Start(ctx); err != nil {
			return err
		}
	}

	go s.simulationLoop(ctx)
	return nil
}

// Stop terminates the simulation
func (s *Simulator) Stop() error {
	close(s.done)
	if err := s.transport.Stop(); err != nil {
		return err
	}
	if s.webSocket != nil {
		return s.webSocket.Stop()
	}
	return nil
}

func (s *Simulator) simulationLoop(ctx context.Context) {
	ticker := time.NewTicker(s.updatePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		case <-ticker.C:
			s.generateAndSendMessages()
		}
	}
}

func (s *Simulator) generateAndSendMessages() {
	// Generate and send vessel heading
	heading := pgn.VesselHeading{
		Heading:   0.5, // ~28.6 degrees
		Reference: 0,   // True heading
	}
	msg := pgn.Message{
		PGN:  127250,
		Data: pgn.EncodeVesselHeading(heading),
	}
	s.transport.SendPGN(msg)
	if s.webSocket != nil {
		s.webSocket.SendPGN(msg)
	}

	// Generate and send speed data
	speed := pgn.SpeedData{
		SpeedWater:  2.5, // 2.5 m/s through water
		SpeedGround: 2.7, // 2.7 m/s over ground
		Reference:   0,   // Paddle wheel
	}
	msg = pgn.Message{
		PGN:  128259,
		Data: pgn.EncodeSpeedData(speed),
	}
	s.transport.SendPGN(msg)
	if s.webSocket != nil {
		s.webSocket.SendPGN(msg)
	}

	// Generate and send water depth
	depth := pgn.WaterDepth{
		Depth:    10.5, // 10.5 meters
		Offset:   -1.5, // Transducer is 1.5m below water line
		MaxRange: 100.0,
	}
	msg = pgn.Message{
		PGN:  128267,
		Data: pgn.EncodeWaterDepth(depth),
	}
	s.transport.SendPGN(msg)
	if s.webSocket != nil {
		s.webSocket.SendPGN(msg)
	}

	// Generate and send position data
	pos := pgn.Position{
		Latitude:  48.1964, // 48°11'47"N
		Longitude: 16.3637, // 16°21'49"E
	}
	msg = pgn.Message{
		PGN:  129025,
		Data: pgn.EncodePosition(pos),
	}
	s.transport.SendPGN(msg)
	if s.webSocket != nil {
		s.webSocket.SendPGN(msg)
	}

	// Generate and send wind data
	wind := pgn.WindData{
		WindSpeed: 5.5, // 5.5 m/s
		WindAngle: 2.0, // ~114.6 degrees
		Reference: 1,   // Apparent wind
	}
	msg = pgn.Message{
		PGN:  130306,
		Data: pgn.EncodeWindData(wind),
	}
	s.transport.SendPGN(msg)
	if s.webSocket != nil {
		s.webSocket.SendPGN(msg)
	}
}
