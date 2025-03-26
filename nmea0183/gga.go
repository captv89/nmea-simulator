// Package nmea0183 provides functions to generate NMEA 0183 sentences.
package nmea0183

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateGGAData generates a GGA sentence.
func GenerateGGAData() string {
	// 1. Time (UTC)
	now := time.Now().UTC()
	utcTime := now.Format("150405.00") // hhmmss.ss

	// 2. Latitude (llll.ll) and N/S
	// For simplicity, let's use a fixed location for now.
	latitude := 4811.7646
	latDirection := "N"

	// 3. Longitude (yyyyy.yy) and E/W
	longitude := 1621.4916
	lonDirection := "E"

	// 4. GPS Quality Indicator
	// 0 = no fix, 1 = GPS fix, 2 = DGPS fix
	gpsQuality := rand.Intn(3)

	// 5. Number of satellites in view
	numSatellites := rand.Intn(12) // Up to 12 satellites

	// 6. Horizontal Dilution of Precision (HDOP)
	hdop := 0.5 + rand.Float64()*4.5 // Between 0.5 and 5.0

	// 7. Antenna Altitude above/below mean-sea-level (geoid)
	altitude := 100.0 + rand.Float64()*200.0 // Example range

	// 8. Geoidal separation
	geoidalSeparation := -5.0 + rand.Float64()*10.0 // Example range

	// 9. Age of differential GPS data (null if not used)
	ageDGPS := ""

	// 10. Differential reference station ID (null if not used)
	refStationID := ""

	// Format the GGA sentence
	ggaSentence := fmt.Sprintf(
		"$GPGGA,%s,%.4f,%s,%.4f,%s,%d,%02d,%.1f,%.1f,M,%.1f,M,%s,%s",
		utcTime,
		latitude, latDirection,
		longitude, lonDirection,
		gpsQuality, numSatellites, hdop,
		altitude, geoidalSeparation,
		ageDGPS, refStationID,
	)

	// Calculate the checksum
	checksum := calculateChecksum(ggaSentence)

	// Append the checksum to the GGA sentence
	ggaSentenceWithChecksum := fmt.Sprintf("%s*%s", ggaSentence, checksum)

	return ggaSentenceWithChecksum
}
