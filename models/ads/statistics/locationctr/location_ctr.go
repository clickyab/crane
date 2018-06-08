package locationctr

import (
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/pool"
)

type pools struct {
	key     string
	driver  pool.Driver
	scanner kv.Scanner
}

var locationCTRPools = make(map[int64]pools)

//GetCRPerLocationByKeys try to get creatives statistics per location
func GetCRPerLocationByKeys(domain string, pageID, seatID, creativeID, creativeSize int64) entity.LocationCTR {
	if ok := locationCTRPools[creativeID]; ok.key == "" {
		return nil
	}

	key := entities.GenCRPerLocationPoolKey(
		domain,
		pageID,
		seatID,
		creativeID,
		creativeSize,
	)

	var data entities.CreativesLocationsReport
	d, err := locationCTRPools[creativeID].driver.Fetch(key, data)
	if err != nil {
		return nil
	}

	return d.(entity.LocationCTR)
}

// AddAndGetCreativePerLocation get crPerLocation by key in pool if not found select on db and if not found again inser it
func AddAndGetCreativePerLocation(page entities.PublisherPage, seat entities.Seat, crID, crSize int64) (entity.LocationCTR, error) {
	crPerLocation := GetCRPerLocationByKeys(
		page.PublisherDomain,
		page.ID,
		seat.ID,
		crID,
		crSize,
	)

	if crPerLocation != nil {
		return crPerLocation, nil
	}
	crReport := entities.CreativesLocationsReport{
		PubID:     seat.PublisherID,
		PubDomain: seat.PublisherDomain,
		PubPageID: page.ID,
		URLKey:    page.URLKey,
		CrID:      crID,
		CrSize:    crSize,
		SID:       seat.ID,
	}
	return entities.AddAndGetCreativesLocationsReport(crReport)
}
