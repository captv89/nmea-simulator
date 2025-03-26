package nmea0183

import (
	"strings"
	"testing"
)

func TestGenerateGGAData(t *testing.T) {
	ggaData := GenerateGGAData()

	// Check if the sentence starts with "$GPGGA"
	if !strings.HasPrefix(ggaData, "$GPGGA") {
		t.Errorf("GGA sentence does not start with $GPGGA: %s", ggaData)
	}

	// Check if the sentence contains a checksum
	if !strings.Contains(ggaData, "*") {
		t.Errorf("GGA sentence does not contain a checksum: %s", ggaData)
	}

	// Split the sentence and checksum
	parts := strings.Split(ggaData, "*")
	if len(parts) != 2 {
		t.Errorf("GGA sentence is not properly formatted with checksum: %s", ggaData)
		return
	}

	// Validate the checksum
	sentence := parts[0][1:] // Remove the initial '$'
	expectedChecksum := calculateChecksum("$" + sentence)
	if parts[1] != expectedChecksum {
		t.Errorf("Checksum mismatch. Expected: %s, Got: %s", expectedChecksum, parts[1])
	}

	// Additional checks for sentence structure
	fields := strings.Split(sentence, ",")
	if len(fields) < 15 {
		t.Errorf("GGA sentence has insufficient fields: %s", ggaData)
	}
}
