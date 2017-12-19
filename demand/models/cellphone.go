package models

import (
	"regexp"

	"clickyab.com/crane/demand/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	networks pool.Interface
	carriers pool.Interface
	brands   pool.Interface
	irmci    = regexp.MustCompile("(?i)(IR)?(-)?(MCI|TCI|43270|Mobile Communications Company of Iran)$")
	irancel  = regexp.MustCompile("(?i)(MTN)?(-)?(irancell|mtn|Iran( )?cell Telecommunications Services Company)$")
	rightel  = regexp.MustCompile("(?i)(righ( )?tel(@)?|IRN 20)$") // Some case are like "Rightle | rightel"
)

// getPhoneData try to insert/retrieve brand for phone
func getPhoneData(carrier string) string {
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
func GetCarrierByName(name string) (int64, error) {
	name = getPhoneData(name)
	d := &entities.Carrier{}
	res, err := carriers.Get(name, d)
	if err != nil {
		return 0, err
	}
	return res.(*entities.Carrier).ID, nil
}

// GetNetworkByName is a function to get the network by its name
func GetNetworkByName(name string) (int64, error) {
	d := &entities.Network{}
	res, err := networks.Get(name, d)
	if err != nil {
		return 0, err
	}
	return res.(*entities.Network).ID, nil
}

// GetBrandByName is a function to get the brand by its name
func GetBrandByName(name string) (int64, error) {
	d := &entities.Brand{}
	res, err := brands.Get(name, d)
	if err != nil {
		return 0, err
	}
	return res.(*entities.Brand).ID, nil
}
