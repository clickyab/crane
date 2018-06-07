package seats

import (
	"strconv"

	"clickyab.com/crane/models/internal/entities"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/pool"
)

var seatsPool pool.Interface

//GetSeats get total seats of all creatives in network per type
//TODO: should use seat interface instead of structure
func GetSeats() []entities.Seat {
	data := seatsPool.All()
	all := make([]entities.Seat, len(data))

	var c int
	for i := range data {
		all[c] = data[i].(entities.Seat)
		c++
	}

	return all
}

// GetSeatByKeys try to get network seats based on its creative type
//TODO: should use seat interface instead of structure
func GetSeatByKeys(supplierName string, slID int64, publisherID int64, publisherDomain string, crSize int64) *entities.Seat {
	data := seatsPool.All()

	key := entities.GenSeatPoolKey(
		supplierName,
		slID,
		publisherID,
		publisherDomain,
		crSize,
	)

	seat := data[key]
	if seat == nil {
		return nil
	}

	return seat.(*entities.Seat)
}

// AddAndGetSeat get seat by key in pool if not found select on db and if not found again inser it
func AddAndGetSeat(m models.Impression, s models.Seat) (*entities.Seat, error) {
	size := int64(s.AdSize)
	sl, _ := strconv.Atoi(s.SlotPublicID)
	slID := int64(sl)

	seat := GetSeatByKeys(
		m.Supplier,
		slID,
		m.PublisherID,
		m.Publisher,
		size,
	)

	if seat != nil {
		return seat, nil
	}

	return entities.AddAndGetSeat(m, size, slID)
}