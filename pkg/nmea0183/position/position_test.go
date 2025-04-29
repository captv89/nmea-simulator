package position

import (
	"strings"
	"testing"
)

func TestGenerateGGA(t *testing.T) {
	gga := GenerateGGA()

	if !strings.HasPrefix(gga, "$GPGGA") {
		t.Errorf("GGA sentence does not start with $GPGGA: %s", gga)
	}

	if !strings.Contains(gga, "*") {
		t.Errorf("GGA sentence does not contain a checksum: %s", gga)
	}

	parts := strings.Split(strings.Split(gga, "*")[0], ",")
	if len(parts) < 15 {
		t.Errorf("GGA sentence has insufficient fields: %s", gga)
	}
}

func TestGenerateGLL(t *testing.T) {
	gll := GenerateGLL()

	if !strings.HasPrefix(gll, "$GPGLL") {
		t.Errorf("GLL sentence should start with $GPGLL, got: %s", gll)
	}

	parts := strings.Split(strings.Split(gll, "*")[0], ",")
	if len(parts) != 7 {
		t.Errorf("Expected 7 fields in GLL sentence, got %d", len(parts))
	}

	timeStr := parts[5]
	if len(timeStr) != 9 { // "HHMMSS.ss" = 9 chars
		t.Errorf("Invalid time format in GLL sentence: %s", timeStr)
	}

	status := parts[6]
	if status != "A" && status != "V" {
		t.Errorf("Invalid status in GLL sentence: %s", status)
	}
}