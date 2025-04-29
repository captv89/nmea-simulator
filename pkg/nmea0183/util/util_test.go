package util

import (
	"strings"
	"testing"
	"time"
)

func TestFormatUTCTime(t *testing.T) {
	testTime := time.Date(2025, 4, 29, 15, 0o4, 0o5, 0, time.UTC)
	result := FormatUTCTime(testTime)

	if len(result) != 9 {
		t.Errorf("Expected time string of length 9, got length %d", len(result))
	}
	if !strings.Contains(result, ".") {
		t.Error("Expected time string to contain decimal point")
	}
	if result != "150405.00" {
		t.Errorf("Expected 150405.00, got %s", result)
	}
}

func TestRandomFloat(t *testing.T) {
	min, max := 0.0, 10.0
	for i := 0; i < 1000; i++ {
		result := RandomFloat(min, max)
		if result < min || result > max {
			t.Errorf("RandomFloat(%f, %f) returned %f, which is outside range", min, max, result)
		}
	}
}

func TestRandomInt(t *testing.T) {
	min, max := 1, 10
	for i := 0; i < 1000; i++ {
		result := RandomInt(min, max)
		if result < min || result > max {
			t.Errorf("RandomInt(%d, %d) returned %d, which is outside range", min, max, result)
		}
	}
}

func TestAppendChecksum(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "GGA example",
			input:    "$GPGGA,092750.000,5321.6802,N,00630.3372,W,1,8,1.03,61.7,M,55.2,M,,",
			expected: "$GPGGA,092750.000,5321.6802,N,00630.3372,W,1,8,1.03,61.7,M,55.2,M,,*76",
		},
		{
			name:     "Simple sentence",
			input:    "$GPGLL,5321.6802,N,00630.3372,W,092750.000,A",
			expected: "$GPGLL,5321.6802,N,00630.3372,W,092750.000,A*26",
		},
		{
			name:     "Empty fields",
			input:    "$GPGGA,,,,,,,,,,,,,,,",
			expected: "$GPGGA,,,,,,,,,,,,,,,*7A",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := AppendChecksum(tc.input)
			if result != tc.expected {
				t.Errorf("AppendChecksum(%s) = %s; want %s", tc.input, result, tc.expected)
			}
		})
	}
}

func TestRandomFloatDistribution(t *testing.T) {
	min, max := 0.0, 1.0
	buckets := make([]int, 10)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		value := RandomFloat(min, max)
		bucketIndex := int(value * 10)
		if bucketIndex == 10 {
			bucketIndex = 9 // Handle edge case of value = 1.0
		}
		buckets[bucketIndex]++
	}

	// Check if distribution is roughly uniform
	expectedCount := iterations / 10
	tolerance := float64(expectedCount) * 0.2 // Allow 20% deviation

	for i, count := range buckets {
		if float64(count) < float64(expectedCount)-tolerance || float64(count) > float64(expectedCount)+tolerance {
			t.Errorf("Bucket %d has count %d, expected roughly %d (±%f)", i, count, expectedCount, tolerance)
		}
	}
}

func TestRandomIntDistribution(t *testing.T) {
	min, max := 1, 6 // Simulate dice rolls
	buckets := make([]int, max-min+1)
	iterations := 6000

	for i := 0; i < iterations; i++ {
		value := RandomInt(min, max)
		buckets[value-min]++
	}

	// Check if distribution is roughly uniform
	expectedCount := iterations / (max - min + 1)
	tolerance := float64(expectedCount) * 0.2 // Allow 20% deviation

	for i, count := range buckets {
		if float64(count) < float64(expectedCount)-tolerance || float64(count) > float64(expectedCount)+tolerance {
			t.Errorf("Value %d has count %d, expected roughly %d (±%f)", i+min, count, expectedCount, tolerance)
		}
	}
}
