package models

import (
	"fmt"
	"net"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/models/internal/entities"
	"clickyab.com/crane/demand/workers/models"
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

// FindPublisher return publisher id for given supplier,domain
func FindPublisher(sup, domain string, pid int64) (entity.Publisher, error) {
	osup, err := GetSupplierByName(sup)
	if err != nil {
		return nil, err
	}
	p, err := GetWebSiteID(osup, domain, pid)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// AddImpression insert new impression to daily table
// TODO : multiple insert per query
func AddImpression(publisher entity.Publisher, impression models.Impression, seat models.Seat) error {
	return entities.AddImpression(publisher, impression, seat)
}

// AdClick try to add new click
func AdClick(supplier, reservedHash, publisher, slotPublicID, referrer, parentURL, os, copID string,
	susp, size int, fast, adID int64, winnerBid float64, ip net.IP, ts time.Time) error {
	// find publisher id
	pub, err := FindPublisher(supplier, publisher, 0)
	if err != nil {
		return err
	}

	click, err := entities.FillClickData(pub.Supplier().Name(), reservedHash, slotPublicID, referrer, parentURL, os, copID,
		susp, size, fast, adID, winnerBid, ip, ts, pub.ID())
	if err != nil {
		return err
	}
	return entities.InsertClick(click)
}
