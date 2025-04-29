// Package navigation provides NMEA-0183 navigation-related sentence generators
package navigation

import (
	"fmt"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/util"
)

// GenerateRMC generates an RMC (Recommended Minimum Navigation Information) sentence
func GenerateRMC() string {
	now := time.Now().UTC()
	utcTime := util.FormatUTCTime(now)
	date := now.Format("020106") // ddmmyy

	status := "A"
	latitude := 4811.7646
	latDirection := "N"
	longitude := 1621.4916
	lonDirection := "E"

	speedKnots := util.RandomFloat(0, 20.0)
	trackTrue := util.RandomFloat(0, 360.0)
	magVar := 5.0 + util.RandomFloat(0, 2.0)
	magVarDirection := "E"

	sentence := fmt.Sprintf(
		"$GPRMC,%s,%s,%.4f,%s,%.4f,%s,%.1f,%.1f,%s,%.1f,%s",
		utcTime, status,
		latitude, latDirection,
		longitude, lonDirection,
		speedKnots, trackTrue,
		date, magVar, magVarDirection,
	)

	return util.AppendChecksum(sentence)
}

// GenerateHDT generates an HDT (Heading - True) sentence
func GenerateHDT() string {
	heading := util.RandomFloat(0, 360.0)

	sentence := fmt.Sprintf(
		"$HEHDT,%.1f,T",
		heading,
	)

	return util.AppendChecksum(sentence)
}

// GenerateVTG generates a VTG (Track Made Good and Ground Speed) sentence
func GenerateVTG() string {
	trackTrue := util.RandomFloat(0, 360.0)
	trackMag := trackTrue - 2.0
	speedKnots := util.RandomFloat(0, 20.0)
	speedKmh := speedKnots * 1.852

	sentence := fmt.Sprintf(
		"$GPVTG,%.1f,T,%.1f,M,%.1f,N,%.1f,K",
		trackTrue, trackMag, speedKnots, speedKmh,
	)

	return util.AppendChecksum(sentence)
}

// GenerateXTE generates an XTE (Cross-Track Error) sentence
func GenerateXTE() string {
	status := "A"
	cycleLock := "A"
	xteDistance := util.RandomFloat(0, 0.5)
	direction := "R"
	if util.RandomFloat(0, 1) < 0.5 {
		direction = "L"
	}
	units := "N"

	sentence := fmt.Sprintf(
		"$GPXTE,%s,%s,%.3f,%s,%s",
		status, cycleLock, xteDistance, direction, units,
	)

	return util.AppendChecksum(sentence)
}
