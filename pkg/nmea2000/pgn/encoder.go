package pgn

import (
	"encoding/binary"
)

// VesselHeading represents PGN 127250 data
type VesselHeading struct {
	Heading   float64 // Radians
	Deviation float64 // Radians
	Variation float64 // Radians
	Reference uint8   // 0=True, 1=Magnetic
	Reserved  uint8
}

// EncodeVesselHeading encodes PGN 127250 data
func EncodeVesselHeading(h VesselHeading) []byte {
	data := make([]byte, 8)

	// Convert heading to radians * 10000
	heading := uint16(h.Heading * 10000)
	binary.LittleEndian.PutUint16(data[0:2], heading)

	// Deviation
	deviation := int16(h.Deviation * 10000)
	binary.LittleEndian.PutUint16(data[2:4], uint16(deviation))

	// Variation
	variation := int16(h.Variation * 10000)
	binary.LittleEndian.PutUint16(data[4:6], uint16(variation))

	// Reference and reserved
	data[6] = h.Reference
	data[7] = h.Reserved

	return data
}

// WaterDepth represents PGN 128267 data
type WaterDepth struct {
	Depth    float64 // Meters
	Offset   float64 // Meters
	MaxRange float64 // Meters
}

// EncodeWaterDepth encodes PGN 128267 data
func EncodeWaterDepth(d WaterDepth) []byte {
	data := make([]byte, 8)

	// Convert depth to centimeters (0.01m resolution)
	depth := uint32(d.Depth * 100)
	binary.LittleEndian.PutUint32(data[0:4], depth)

	// Offset
	offset := int16(d.Offset * 100)
	binary.LittleEndian.PutUint16(data[4:6], uint16(offset))

	// Range
	maxRange := uint16(d.MaxRange * 100)
	binary.LittleEndian.PutUint16(data[6:8], maxRange)

	return data
}

// WindData represents PGN 130306 data
type WindData struct {
	WindSpeed float64 // Meters per second
	WindAngle float64 // Radians
	Reference uint8   // 0=True, 1=Apparent
}

// EncodeWindData encodes PGN 130306 data
func EncodeWindData(w WindData) []byte {
	data := make([]byte, 8)

	// Wind speed (0.01 m/s resolution)
	speed := uint16(w.WindSpeed * 100)
	binary.LittleEndian.PutUint16(data[0:2], speed)

	// Wind angle (0.0001 radian resolution)
	angle := uint16(w.WindAngle * 10000)
	binary.LittleEndian.PutUint16(data[2:4], angle)

	// Reference
	data[4] = w.Reference

	// Reserved bytes
	data[5] = 0xFF
	data[6] = 0xFF
	data[7] = 0xFF

	return data
}

// Position represents PGN 129025 data
type Position struct {
	Latitude  float64 // Degrees
	Longitude float64 // Degrees
}

// EncodePosition encodes PGN 129025 data
func EncodePosition(p Position) []byte {
	data := make([]byte, 8)

	// Convert latitude to 1e-7 degrees
	lat := int32(p.Latitude * 1e7)
	binary.LittleEndian.PutUint32(data[0:4], uint32(lat))

	// Convert longitude to 1e-7 degrees
	lon := int32(p.Longitude * 1e7)
	binary.LittleEndian.PutUint32(data[4:8], uint32(lon))

	return data
}

// SpeedData represents PGN 128259 data
type SpeedData struct {
	SpeedWater  float64 // Meters per second
	SpeedGround float64 // Meters per second
	Reference   uint8   // 0=Paddle wheel, 1=Pitot tube, 2=Doppler, 3=Correlation
}

// EncodeSpeedData encodes PGN 128259 data
func EncodeSpeedData(s SpeedData) []byte {
	data := make([]byte, 8)

	// Water speed (0.01 m/s resolution)
	waterSpeed := uint16(s.SpeedWater * 100)
	binary.LittleEndian.PutUint16(data[0:2], waterSpeed)

	// Ground speed
	groundSpeed := uint16(s.SpeedGround * 100)
	binary.LittleEndian.PutUint16(data[2:4], groundSpeed)

	// Reference type
	data[4] = s.Reference

	// Reserved bytes
	data[5] = 0xFF
	data[6] = 0xFF
	data[7] = 0xFF

	return data
}
