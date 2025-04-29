// Package util provides utility functions for NMEA sentence generation
package util

import (
	"fmt"
	"math/rand"
	"time"
)

// FormatUTCTime formats time in NMEA UTC format (hhmmss.ss)
func FormatUTCTime(t time.Time) string {
	return t.Format("150405.00")
}

// RandomFloat generates a random float64 between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomInt generates a random integer between min and max (inclusive)
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// AppendChecksum calculates and appends the checksum to an NMEA sentence.
// The checksum is calculated by XOR'ing all characters between $ and * (exclusive).
func AppendChecksum(sentence string) string {
	var checksum uint8
	var i int

	// Find the start of the sentence (after $)
	for i = 0; i < len(sentence); i++ {
		if sentence[i] == '$' {
			i++
			break
		}
	}

	// Calculate checksum from character after $ until * or end of string
	for ; i < len(sentence) && sentence[i] != '*'; i++ {
		checksum ^= sentence[i]
	}

	// Format the checksum as a two-character uppercase hexadecimal
	return fmt.Sprintf("%s*%02X", sentence, checksum)
}
