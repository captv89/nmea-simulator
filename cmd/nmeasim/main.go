package main

import (
	"fmt"
	"time"

	"github.com/captv89/nmea-simulator/pkg/nmea0183/environment"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/navigation"
	"github.com/captv89/nmea-simulator/pkg/nmea0183/position"
)

func main() {
	// Run continuously with a 1-second interval
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Position sentences
		fmt.Println(position.GenerateGGA())
		fmt.Println(position.GenerateGLL())

		// Navigation sentences
		fmt.Println(navigation.GenerateRMC())
		fmt.Println(navigation.GenerateHDT())
		fmt.Println(navigation.GenerateVTG())
		fmt.Println(navigation.GenerateXTE())

		// Environment sentences
		fmt.Println(environment.GenerateDBT())
		fmt.Println(environment.GenerateMTW())
		fmt.Println(environment.GenerateMWV())
		fmt.Println(environment.GenerateVHW())
		fmt.Println(environment.GenerateDPT())

		fmt.Println("---") // Separator between sets
	}
}
