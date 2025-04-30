package pgn

import (
	"encoding/binary"
	"math"
	"testing"
)

func float64ToUint16(f float64) uint16 {
	if f < 0 {
		i := int16(f)
		return uint16(i)
	}
	return uint16(f)
}

func float64ToUint32(f float64) uint32 {
	if f < 0 {
		i := int32(f)
		return uint32(i)
	}
	return uint32(f)
}

func TestEncodeVesselHeading(t *testing.T) {
	tests := []struct {
		name     string
		heading  VesselHeading
		wantData []byte
	}{
		{
			name: "true heading with zero deviation and variation",
			heading: VesselHeading{
				Heading:   0.5, // ~28.6 degrees
				Deviation: 0,
				Variation: 0,
				Reference: 0, // True heading
				Reserved:  0,
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16(0.5*10000))
				return data
			}(),
		},
		{
			name: "magnetic heading with deviation and variation",
			heading: VesselHeading{
				Heading:   math.Pi / 2, // 90 degrees
				Deviation: 0.1,         // ~5.7 degrees
				Variation: -0.05,       // ~-2.9 degrees
				Reference: 1,           // Magnetic heading
				Reserved:  0,
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16((math.Pi/2)*10000))
				binary.LittleEndian.PutUint16(data[2:4], float64ToUint16(0.1*10000))
				binary.LittleEndian.PutUint16(data[4:6], float64ToUint16(-0.05*10000))
				data[6] = 1 // Reference
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeVesselHeading(tt.heading)
			if len(got) != 8 {
				t.Errorf("EncodeVesselHeading() returned %d bytes, want 8", len(got))
			}
			for i := 0; i < 8; i++ {
				if got[i] != tt.wantData[i] {
					t.Errorf("EncodeVesselHeading() byte[%d] = %v, want %v", i, got[i], tt.wantData[i])
				}
			}
		})
	}
}

func TestEncodeWaterDepth(t *testing.T) {
	tests := []struct {
		name     string
		depth    WaterDepth
		wantData []byte
	}{
		{
			name: "typical depth with offset",
			depth: WaterDepth{
				Depth:    10.5, // 10.5 meters
				Offset:   -1.5, // Transducer is 1.5m below water line
				MaxRange: 100.0,
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint32(data[0:4], float64ToUint32(10.5*100))
				binary.LittleEndian.PutUint16(data[4:6], float64ToUint16(-1.5*100))
				binary.LittleEndian.PutUint16(data[6:8], float64ToUint16(100.0*100))
				return data
			}(),
		},
		{
			name: "shallow depth no offset",
			depth: WaterDepth{
				Depth:    2.0,
				Offset:   0,
				MaxRange: 50.0,
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint32(data[0:4], uint32(2.0*100))
				binary.LittleEndian.PutUint16(data[6:8], uint16(50.0*100))
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeWaterDepth(tt.depth)
			if len(got) != 8 {
				t.Errorf("EncodeWaterDepth() returned %d bytes, want 8", len(got))
			}
			for i := 0; i < 8; i++ {
				if got[i] != tt.wantData[i] {
					t.Errorf("EncodeWaterDepth() byte[%d] = %v, want %v", i, got[i], tt.wantData[i])
				}
			}
		})
	}
}

func TestEncodeWindData(t *testing.T) {
	tests := []struct {
		name     string
		wind     WindData
		wantData []byte
	}{
		{
			name: "apparent wind",
			wind: WindData{
				WindSpeed: 5.5, // 5.5 m/s
				WindAngle: 2.0, // ~114.6 degrees
				Reference: 1,   // Apparent wind
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16(5.5*100))
				binary.LittleEndian.PutUint16(data[2:4], float64ToUint16(2.0*10000))
				data[4] = 1 // Reference
				data[5] = 0xFF
				data[6] = 0xFF
				data[7] = 0xFF
				return data
			}(),
		},
		{
			name: "true wind",
			wind: WindData{
				WindSpeed: 10.0,
				WindAngle: math.Pi, // 180 degrees
				Reference: 0,       // True wind
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16(10.0*100))
				binary.LittleEndian.PutUint16(data[2:4], float64ToUint16(math.Pi*10000))
				data[5] = 0xFF
				data[6] = 0xFF
				data[7] = 0xFF
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeWindData(tt.wind)
			if len(got) != 8 {
				t.Errorf("EncodeWindData() returned %d bytes, want 8", len(got))
			}
			for i := 0; i < 8; i++ {
				if got[i] != tt.wantData[i] {
					t.Errorf("EncodeWindData() byte[%d] = %v, want %v", i, got[i], tt.wantData[i])
				}
			}
		})
	}
}

func TestEncodePosition(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		wantData []byte
	}{
		{
			name: "positive lat/lon",
			pos: Position{
				Latitude:  48.1964, // 48°11'47"N
				Longitude: 16.3637, // 16°21'49"E
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint32(data[0:4], float64ToUint32(48.1964*1e7))
				binary.LittleEndian.PutUint32(data[4:8], float64ToUint32(16.3637*1e7))
				return data
			}(),
		},
		{
			name: "negative lat/lon",
			pos: Position{
				Latitude:  -33.8688,  // 33°52'08"S
				Longitude: -151.2093, // 151°12'33"W
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint32(data[0:4], float64ToUint32(-33.8688*1e7))
				binary.LittleEndian.PutUint32(data[4:8], float64ToUint32(-151.2093*1e7))
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodePosition(tt.pos)
			if len(got) != 8 {
				t.Errorf("EncodePosition() returned %d bytes, want 8", len(got))
			}
			for i := 0; i < 8; i++ {
				if got[i] != tt.wantData[i] {
					t.Errorf("EncodePosition() byte[%d] = %v, want %v", i, got[i], tt.wantData[i])
				}
			}
		})
	}
}

func TestEncodeSpeedData(t *testing.T) {
	tests := []struct {
		name     string
		speed    SpeedData
		wantData []byte
	}{
		{
			name: "paddle wheel speed",
			speed: SpeedData{
				SpeedWater:  2.5, // 2.5 m/s
				SpeedGround: 2.7, // 2.7 m/s
				Reference:   0,   // Paddle wheel
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16(2.5*100))
				binary.LittleEndian.PutUint16(data[2:4], float64ToUint16(2.7*100))
				data[5] = 0xFF
				data[6] = 0xFF
				data[7] = 0xFF
				return data
			}(),
		},
		{
			name: "doppler speed",
			speed: SpeedData{
				SpeedWater:  5.0,
				SpeedGround: 5.2,
				Reference:   2, // Doppler
			},
			wantData: func() []byte {
				data := make([]byte, 8)
				binary.LittleEndian.PutUint16(data[0:2], float64ToUint16(5.0*100))
				binary.LittleEndian.PutUint16(data[2:4], float64ToUint16(5.2*100))
				data[4] = 2
				data[5] = 0xFF
				data[6] = 0xFF
				data[7] = 0xFF
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeSpeedData(tt.speed)
			if len(got) != 8 {
				t.Errorf("EncodeSpeedData() returned %d bytes, want 8", len(got))
			}
			for i := 0; i < 8; i++ {
				if got[i] != tt.wantData[i] {
					t.Errorf("EncodeSpeedData() byte[%d] = %v, want %v", i, got[i], tt.wantData[i])
				}
			}
		})
	}
}
