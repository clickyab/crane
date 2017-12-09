package models

import (
	"fmt"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var ads pool.Interface

// GetAds return all ads in system
func GetAds() []entity.Advertise {
	data := ads.All()
	all := make([]entity.Advertise, len(data))
	var c int
	for i := range data {
		all[c] = data[i].(entity.Advertise)
		c++
	}

	return all
}

// GetAd try to get advertise based on its id
func GetAd(adID int64) (entity.Advertise, error) {
	ad, err := ads.Get(fmt.Sprint(adID), &entities.Advertise{})
	if err != nil {
		ad, err = entities.GetAd(adID)
		if err != nil {
			return nil, err
		}
	}
	return ad.(entity.Advertise), nil
}
