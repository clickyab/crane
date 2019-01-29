package campaign

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/internal/entities"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/pool"
	"github.com/sirupsen/logrus"
)

var campaign pool.Interface

// GetCampaigns return all ads in system
func GetCampaigns() map[int32]entity.Campaign {
	data := campaign.All()
	all := make(map[int32]entity.Campaign)
	for i := range data {
		c := data[i].(entity.Campaign)
		all[c.ID()] = c
	}
	return all
}

// GetAd try to get advertise based on its id
func GetAd(adID int32) (entity.Creative, error) {
	ad, err := campaign.Get(fmt.Sprint(adID), &entities.Advertise{})
	if err != nil {
		x, err := entities.GetAd(adID)
		if err != nil {
			return nil, err
		}
		return x, nil
	}
	return ad.(entity.Creative), nil
}

// FindPublisher return publisher id for given supplier,domain
func FindPublisher(sup, domain string, pid int64, t entity.PublisherType) (entity.Publisher, error) {
	osup, err := suppliers.GetSupplierByName(sup)
	if err != nil {
		return nil, err
	}
	var p entity.Publisher
	if t == entity.PublisherTypeWeb {
		p, err = website.GetWebSite(osup, domain, pid)
	} else if t == entity.PublisherTypeApp {
		p, err = apps.GetApp(osup, domain, "")
	} else {
		logrus.Errorf("[BUG] invalid request type the type is %s", t.String())
		return nil, fmt.Errorf("[BUG] invalid request type the type is %s", t.String())
	}
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

// AdClick will get impression from job and insert it into click table
func AdClick(p entity.Publisher, m models.Impression, s models.Seat,
	os entity.OS, fast int64, tv bool) error {
	click, err := entities.FillClickData(p, m, s, os, fast, tv)
	if err != nil {
		return err
	}
	return entities.InsertClick(click)
}

// AddNotice get impression from job abd insert it into notice table
func AddNotice(p entity.Publisher, m models.Impression, s models.Seat) error {
	return entities.AddNotice(p, m, s)
}
