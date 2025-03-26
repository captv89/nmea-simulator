package nmea0183

import "fmt"

// calculateChecksum calculates the NMEA 0183 checksum for a sentence.
func calculateChecksum(sentence string) string {
	checksum := 0
	for i := 1; i < len(sentence); i++ {
		checksum ^= int(sentence[i])
	}
	return fmt.Sprintf("%02X", checksum)
}
