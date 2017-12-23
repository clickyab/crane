package models

import (
	"regexp"
)

var (
	irmci   = regexp.MustCompile("(?i)(IR)?(-)?(MCI|TCI|43270|Mobile Communications Company of Iran)$")
	irancel = regexp.MustCompile("(?i)(MTN)?(-)?(irancell|mtn|Iran( )?cell Telecommunications Services Company)$")
	rightel = regexp.MustCompile("(?i)(righ( )?tel(@)?|IRN 20)$") // Some case are like "Rightle | rightel"
)

// sanitizeCarrier try to insert/retrieve brand for phone
func sanitizeCarrier(carrier string) string {
	if irancel.MatchString(carrier) {
		carrier = "Irancell"
	} else if irmci.MatchString(carrier) {
		carrier = "IR-MCI"
	} else if rightel.MatchString(carrier) {
		carrier = "RighTel"
	}
	return carrier
}

// GetCarrierByName is a function to get the carrier by its name
func GetCarrierByName(name string) (string, error) {
	name = sanitizeCarrier(name)
	return name, nil
}

// GetNetworkByName is a function to get the network by its name
func GetNetworkByName(name string) (string, error) {
	return name, nil
}

// GetBrandByName is a function to get the brand by its name
func GetBrandByName(name string) (string, error) {
	return name, nil
}
