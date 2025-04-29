// Package position provides NMEA-0183 position-related sentence generators
package position

import (
	"fmt"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/util"
)

// GenerateGGA generates a GGA (Global Positioning System Fix Data) sentence
func GenerateGGA() string {
	// 1. Time (UTC)
	now := time.Now().UTC()
	utcTime := util.FormatUTCTime(now)

	// 2. Latitude and N/S
	latitude := 4811.7646
	latDirection := "N"

	// 3. Longitude and E/W
	longitude := 1621.4916
	lonDirection := "E"

	// 4. GPS Quality Indicator
	gpsQuality := util.RandomInt(0, 2)

	// 5. Number of satellites
	numSatellites := util.RandomInt(4, 12)

	// 6. HDOP
	hdop := 0.5 + util.RandomFloat(0, 4.5)

	// 7. Altitude
	altitude := 100.0 + util.RandomFloat(0, 200.0)

	// 8. Geoidal separation
	geoidalSeparation := -5.0 + util.RandomFloat(0, 10.0)

	// Format the sentence
	sentence := fmt.Sprintf(
		"$GPGGA,%s,%.4f,%s,%.4f,%s,%d,%02d,%.1f,%.1f,M,%.1f,M,,",
		utcTime, latitude, latDirection, longitude, lonDirection,
		gpsQuality, numSatellites, hdop, altitude, geoidalSeparation,
	)

	return util.AppendChecksum(sentence)
}

// GenerateGLL generates a GLL (Geographic Position - Latitude/Longitude) sentence
func GenerateGLL() string {
	// Reuse position data for consistency
	latitude := 4811.7646
	latDirection := "N"
	longitude := 1621.4916
	lonDirection := "E"

	now := time.Now().UTC()
	utcTime := util.FormatUTCTime(now)

	status := "A"

	sentence := fmt.Sprintf(
		"$GPGLL,%.4f,%s,%.4f,%s,%s,%s",
		latitude, latDirection,
		longitude, lonDirection,
		utcTime, status,
	)

	return util.AppendChecksum(sentence)
}
