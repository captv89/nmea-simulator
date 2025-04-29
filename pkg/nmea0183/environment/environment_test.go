package environment

import (
	"strconv"
	"strings"
	"testing"
)

func TestGenerateDBT(t *testing.T) {
	dbt := GenerateDBT()

	if !strings.HasPrefix(dbt, "$IIDBT") {
		t.Errorf("DBT sentence should start with $IIDBT, got: %s", dbt)
	}

	parts := strings.Split(strings.Split(dbt, "*")[0], ",")
	if len(parts) != 7 {
		t.Errorf("Expected 7 fields in DBT sentence, got %d", len(parts))
	}

	if parts[2] != "f" || parts[4] != "M" || parts[6] != "F" {
		t.Error("Invalid units in DBT sentence")
	}
}

func TestGenerateMTW(t *testing.T) {
	mtw := GenerateMTW()

	if !strings.HasPrefix(mtw, "$IIMTW") {
		t.Errorf("MTW sentence should start with $IIMTW, got: %s", mtw)
	}

	parts := strings.Split(strings.Split(mtw, "*")[0], ",")
	if len(parts) != 3 {
		t.Errorf("Expected 3 fields in MTW sentence, got %d", len(parts))
	}

	temp, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || temp < 10.0 || temp > 30.0 {
		t.Errorf("Invalid temperature value in MTW sentence: %s", parts[1])
	}

	if parts[2] != "C" {
		t.Error("Temperature unit should be C")
	}
}

func TestGenerateMWV(t *testing.T) {
	mwv := GenerateMWV()

	if !strings.HasPrefix(mwv, "$IIMWV") {
		t.Errorf("MWV sentence should start with $IIMWV, got: %s", mwv)
	}

	parts := strings.Split(strings.Split(mwv, "*")[0], ",")
	if len(parts) != 6 {
		t.Errorf("Expected 6 fields in MWV sentence, got %d", len(parts))
	}

	angle, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || angle < 0 || angle > 360 {
		t.Errorf("Invalid wind angle in MWV sentence: %s", parts[1])
	}

	if parts[2] != "R" || parts[4] != "N" || parts[5] != "A" {
		t.Error("Invalid reference, speed unit or status in MWV sentence")
	}
}

func TestGenerateDPT(t *testing.T) {
	dpt := GenerateDPT()

	if !strings.HasPrefix(dpt, "$IIDPT") {
		t.Errorf("DPT sentence should start with $IIDPT, got: %s", dpt)
	}

	parts := strings.Split(strings.Split(dpt, "*")[0], ",")
	if len(parts) != 4 {
		t.Errorf("Expected 4 fields in DPT sentence, got %d", len(parts))
	}

	depth, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || depth < 5.0 || depth > 100.0 {
		t.Error("Invalid depth value in DPT sentence")
	}
}

func TestGenerateVHW(t *testing.T) {
	vhw := GenerateVHW()

	if !strings.HasPrefix(vhw, "$IIVHW") {
		t.Errorf("VHW sentence should start with $IIVHW, got: %s", vhw)
	}

	parts := strings.Split(strings.Split(vhw, "*")[0], ",")
	if len(parts) != 9 {
		t.Errorf("Expected 9 fields in VHW sentence, got %d", len(parts))
	}

	if parts[2] != "T" || parts[4] != "M" || parts[6] != "N" || parts[8] != "K" {
		t.Error("Invalid units in VHW sentence")
	}
}