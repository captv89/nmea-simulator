// Package main is the entry point for the NMEA simulator application.
package main

import (
	"fmt"

	"github.com/captv89/nmea-simulator/nmea0183"
)

func main() {
	ggaData := nmea0183.GenerateGGAData()

	fmt.Println(ggaData)
}
