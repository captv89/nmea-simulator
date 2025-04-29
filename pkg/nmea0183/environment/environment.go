// Package environment provides NMEA-0183 environment-related sentence generators
package environment

import (
	"fmt"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/util"
)

// GenerateDBT generates a DBT (Depth Below Transducer) sentence
func GenerateDBT() string {
	depthMeters := 5.0 + util.RandomFloat(0, 95.0)
	depthFeet := depthMeters * 3.28084
	depthFathoms := depthMeters * 0.546807

	sentence := fmt.Sprintf(
		"$IIDBT,%.1f,f,%.1f,M,%.1f,F",
		depthFeet, depthMeters, depthFathoms,
	)

	return util.AppendChecksum(sentence)
}

// GenerateMTW generates an MTW (Mean Temperature of Water) sentence
func GenerateMTW() string {
	tempC := 10.0 + util.RandomFloat(0, 20.0)

	sentence := fmt.Sprintf(
		"$IIMTW,%.1f,C",
		tempC,
	)

	return util.AppendChecksum(sentence)
}

// GenerateMWV generates an MWV (Wind Speed and Angle) sentence
func GenerateMWV() string {
	windAngle := util.RandomFloat(0, 360.0)
	windSpeed := util.RandomFloat(0, 60.0)

	reference := "R"
	speedUnit := "N"
	status := "A"

	sentence := fmt.Sprintf(
		"$IIMWV,%.1f,%s,%.1f,%s,%s",
		windAngle, reference, windSpeed, speedUnit, status,
	)

	return util.AppendChecksum(sentence)
}

// GenerateDPT generates a DPT (Depth of Water) sentence
func GenerateDPT() string {
	depthMeters := 5.0 + util.RandomFloat(0, 95.0)
	offset := -1.5
	maxRange := 200.0

	sentence := fmt.Sprintf(
		"$IIDPT,%.1f,%.1f,%.1f",
		depthMeters, offset, maxRange,
	)

	return util.AppendChecksum(sentence)
}

// GenerateVHW generates a VHW (Water Speed and Heading) sentence
func GenerateVHW() string {
	headingTrue := util.RandomFloat(0, 360.0)
	headingMagnetic := headingTrue - 2.0
	speedKnots := util.RandomFloat(0, 20.0)
	speedKmh := speedKnots * 1.852

	sentence := fmt.Sprintf(
		"$IIVHW,%.1f,T,%.1f,M,%.1f,N,%.1f,K",
		headingTrue, headingMagnetic, speedKnots, speedKmh,
	)

	return util.AppendChecksum(sentence)
}
