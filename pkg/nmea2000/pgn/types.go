// Package pgn provides NMEA 2000 PGN (Parameter Group Number) definitions and encoding
package pgn

// PGNDefinition represents the structure of a NMEA 2000 PGN
type PGNDefinition struct {
	PGN         uint32
	Name        string
	Length      uint8
	Resolution  float64
	Min         float64
	Max         float64
	Units       string
	Description string
}

// CommonPGNs defines the most commonly used PGNs in marine applications
var CommonPGNs = map[uint32]PGNDefinition{
	127250: {
		PGN:         127250,
		Name:        "Vessel Heading",
		Description: "Heading sensor value with a flag for True or Magnetic",
		Length:      8,
	},
	128259: {
		PGN:         128259,
		Name:        "Speed",
		Description: "Speed through water",
		Length:      8,
	},
	128267: {
		PGN:         128267,
		Name:        "Water Depth",
		Description: "Water depth information",
		Length:      8,
	},
	129025: {
		PGN:         129025,
		Name:        "Position Rapid Update",
		Description: "Provides lat/lon rapid update",
		Length:      8,
	},
	129026: {
		PGN:         129026,
		Name:        "COG & SOG Rapid Update",
		Description: "Course Over Ground and Speed Over Ground",
		Length:      8,
	},
	130306: {
		PGN:         130306,
		Name:        "Wind Data",
		Description: "Wind speed, direction, and reference",
		Length:      8,
	},
}

// Message represents a NMEA 2000 message with its PGN and data
type Message struct {
	PGN  uint32
	Data []byte
}