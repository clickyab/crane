package models

import (
	"errors"
	"fmt"
	"net"
	"time"

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
		x, err := entities.GetAd(adID)
		if err != nil {
			return nil, err
		}
		return x, nil
	}
	return ad.(entity.Advertise), nil
}

// ErrorNotAllowCreate rise when supplier doesn't allow to add new website or app
var ErrorNotAllowCreate = errors.New("insert not allowed")

// FindPublisherID return publisher id for given supplier,domain
func FindPublisherID(sup, domain string, pid int64) (int64, error) {
	osup, err := GetSupplierByName(sup)
	if err != nil {
		return 0, err
	}
	p, err := GetWebSiteID(osup, domain, pid)
	if err != nil {
		return 0, err
	}
	return p, nil
}

// AddImpression insert new impression to daily table
func AddImpression(supp, reserve, publisher, ref, parent, spid, copID string, size, susp int, adid int64, ip net.IP,
	bid float64, alexa bool, ts time.Time, typ entity.RequestType) error {

	tw, err := FindPublisherID(supp, publisher, 0)
	if err != nil {
		return err
	}

	return entities.AddImpression(reserve, ref, parent, spid, copID, size, susp, adid, tw, ip, bid, alexa, ts, typ)
}
