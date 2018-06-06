package perlocations

import (
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var crPerLocationsPool pool.Interface

//GetcrPerLocations get total seats of all creatives in network per type
//TODO: should use seat interface instead of structure
func GetcrPerLocations() []entities.CreativesLocationsReport {
	data := crPerLocationsPool.All()
	all := make([]entities.CreativesLocationsReport, len(data))

	var c int
	for i := range data {
		all[c] = data[i].(entities.CreativesLocationsReport)
		c++
	}

	return all
}

//GetCRPerLocationByKeys try to get creatives statistics per location
//TODO: should use seat interface instead of structure
func GetCRPerLocationByKeys(publisherDomain string, pageID, seatID, creativeID int64) *entities.CreativesLocationsReport {
	data := crPerLocationsPool.All()

	key := entities.GenCRPerLocationPoolKey(
		publisherDomain,
		pageID,
		seatID,
		creativeID,
	)

	crPerLocation := data[key]
	if crPerLocation == nil {
		return nil
	}

	return crPerLocation.(*entities.CreativesLocationsReport)
}

// AddAndGetCreativePerLocation get crPerLocation by key in pool if not found select on db and if not found again inser it
func AddAndGetCreativePerLocation(seat entities.Seat, page entities.PublisherPage, crID, crSize int64) (*entities.CreativesLocationsReport, error) {
	crPerLocation := GetCRPerLocationByKeys(
		seat.PublisherDomain,
		page.ID,
		seat.ID,
		crID,
	)

	if crPerLocation != nil {
		return crPerLocation, nil
	}
	crReport := entities.CreativesLocationsReport{
		PublisherID:     seat.PublisherID,
		PublisherDomain: seat.PublisherDomain,
		PublisherPageID: page.ID,
		URLKey:          page.URLKey,
		CreativeID:      crID,
		CreativeSize:    crSize,
		SeatID:          seat.ID,
	}
	return entities.AddAndGetCreativesLocationsReport(crReport)
}
