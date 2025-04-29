package navigation

import (
	"strconv"
	"strings"
	"testing"
)

func TestGenerateRMC(t *testing.T) {
	rmc := GenerateRMC()

	if !strings.HasPrefix(rmc, "$GPRMC") {
		t.Errorf("RMC sentence should start with $GPRMC, got: %s", rmc)
	}

	parts := strings.Split(strings.Split(rmc, "*")[0], ",")
	if len(parts) != 12 {
		t.Errorf("Expected 12 fields in RMC sentence, got %d", len(parts))
	}

	// Check status
	if parts[2] != "A" && parts[2] != "V" {
		t.Errorf("Invalid status in RMC sentence: %s", parts[2])
	}
}

func TestGenerateHDT(t *testing.T) {
	hdt := GenerateHDT()

	if !strings.HasPrefix(hdt, "$HEHDT") {
		t.Errorf("HDT sentence should start with $HEHDT, got: %s", hdt)
	}

	parts := strings.Split(strings.Split(hdt, "*")[0], ",")
	if len(parts) != 3 {
		t.Errorf("Expected 3 fields in HDT sentence, got %d", len(parts))
	}

	heading, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || heading < 0 || heading > 360 {
		t.Errorf("Invalid heading value in HDT sentence: %s", parts[1])
	}
}

func TestGenerateVTG(t *testing.T) {
	vtg := GenerateVTG()

	if !strings.HasPrefix(vtg, "$GPVTG") {
		t.Errorf("VTG sentence should start with $GPVTG, got: %s", vtg)
	}

	parts := strings.Split(strings.Split(vtg, "*")[0], ",")
	if len(parts) != 9 {
		t.Errorf("Expected 9 fields in VTG sentence, got %d", len(parts))
	}

	if parts[2] != "T" || parts[4] != "M" || parts[6] != "N" || parts[8] != "K" {
		t.Error("Invalid units in VTG sentence")
	}
}

func TestGenerateXTE(t *testing.T) {
	xte := GenerateXTE()

	if !strings.HasPrefix(xte, "$GPXTE") {
		t.Errorf("XTE sentence should start with $GPXTE, got: %s", xte)
	}

	parts := strings.Split(strings.Split(xte, "*")[0], ",")
	if len(parts) != 6 {
		t.Errorf("Expected 6 fields in XTE sentence, got %d", len(parts))
	}

	if (parts[1] != "A" && parts[1] != "V") || (parts[2] != "A" && parts[2] != "V") {
		t.Error("Invalid status values in XTE sentence")
	}

	if parts[4] != "L" && parts[4] != "R" {
		t.Errorf("Invalid direction in XTE sentence: %s", parts[4])
	}
}