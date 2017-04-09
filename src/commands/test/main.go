package main

import (
	"commands"
	"fmt"
	"math/rand"
	"services/ip2location"
)

func randInt() int {
	return rand.Intn(254) + 1
}

func main() {
	ip2location.Open()

	for i := 0; i < 100; i++ {

		ip := fmt.Sprintf("%d.%d.%d.%d", randInt(), randInt(), randInt(), randInt())

		results := ip2location.Get_all(ip)

		fmt.Printf("ip: %s\n", ip)
		fmt.Printf("country_short: %s\n", results.Country_short)
		fmt.Printf("country_long: %s\n", results.Country_long)
		fmt.Printf("region: %s\n", results.Region)
		fmt.Printf("city: %s\n", results.City)
		fmt.Printf("latitude: %f\n", results.Latitude)
		fmt.Printf("longitude: %f\n", results.Longitude)
		fmt.Printf("elevation: %f\n", results.Elevation)
	}
	ip2location.Close()
	commands.WaitExitSignal()
}
