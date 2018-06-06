package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/sirupsen/logrus"
)

//TODO: add model codegen
//TODO: fix seat interface and implement it

// Seat seats model in database
// @Model {
//		table = seats
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Seat struct {
	ID              int64                `json:"id" db:"id"`
	SlotID          int64                `json:"slot_id" db:"slot_id"`
	SupplierName    string               `json:"supplier_name" db:"supplier_name"`
	PublisherID     int64                `json:"publisher_id" db:"publisher_id"`
	PublisherDomain string               `json:"publisher_domain" db:"publisher_domain"`
	CreativeSize    int64                `json:"creative_size" db:"creative_size"`
	Kind            entity.PublisherType `json:"kind" db:"kind"`
	ActiveDays      int64                `json:"active_days" db:"active_days"`
	AvgDailyImp     int64                `json:"avg_daily_imp" db:"avg_daily_imp"`
	AvgDailyClicks  int64                `json:"avg_daily_clicks" db:"avg_daily_clicks"`
	TodayImp        int64                `json:"today_imp" db:"today_imp"`
	TodayClicks     int64                `json:"today_clicks" db:"today_clicks"`
	TodayCTR        int64                `json:"today_ctr" db:"today_ctr"`
	CreatedAt       time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt       mysql.NullTime       `json:"updated_at" db:"updated_at"`
}

// AddAndGetSeat return seat if exist and insert if not
func AddAndGetSeat(m models.Impression, crSize, slID int64) (*Seat, error) {
	var seat Seat

	fSeatQ := `SELECT 
			id,
			slot_id,
			supplier_name,
			publisher_id,
			publisher_domain,
			creative_size
		FROM seats
		WHERE
			slot_id=?
			AND supplier_name=?
			AND publisher_id=?
			AND publisher_domain=?
			AND creative_size=?
	`

	//Important: use GetWDbMap because read db may take time to synce and fire err and finally miss impression
	err := NewManager().GetWDbMap().SelectOne(
		&seat,
		fSeatQ,
		slID,
		m.Supplier,
		m.PublisherID,
		m.Publisher,
		crSize,
	)
	if err == nil {
		return &seat, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		seat.SlotID = slID
		seat.SupplierName = m.Supplier
		seat.PublisherID = m.PublisherID
		seat.PublisherDomain = m.Publisher
		seat.CreativeSize = crSize
		seat.Kind = m.PublisherType
		seat.CreatedAt = time.Now()

		err = NewManager().GetWDbMap().Insert(&seat)
		if err != nil {
			logrus.Debug(err)
			return nil, err
		}
	}

	return &seat, nil
}

// SeatsLoader load all seats
func SeatsLoader() func(ctx context.Context) (map[string]kv.Serializable, error) {
	return func(ctx context.Context) (map[string]kv.Serializable, error) {
		seatData := make(map[string]kv.Serializable)
		// return b, nil // Uncomment this line after first time in DEV mode

		const cnt = 10000
		for j := 0; ; j = j + cnt {
			q := fmt.Sprintf(`SELECT 
					id,
					slot_id,
					supplier_name,
					publisher_id,
					publisher_domain,
					creative_size,
					kind,
					active_days,
					today_imp,
					today_clicks,
					today_ctr
				FROM seats
				WHERE 1
				LIMIT %d, %d`,
				j,
				j+cnt,
			)

			var res []Seat
			if _, err := NewManager().GetRDbMap().Select(&res, q); err != nil {
				logrus.Warn(err)
				return nil, err
			}

			if len(res) == 0 {
				break
			}

			for i := range res {
				sID := GenSeatPoolKey(
					res[i].SupplierName,
					res[i].SlotID,
					res[i].PublisherID,
					res[i].PublisherDomain,
					res[i].CreativeSize,
				)
				seatData[sID] = &res[i]
			}
		}

		logrus.Debugf("Load %d seats", len(seatData))

		return seatData, nil
	}
}

// GenSeatPoolKey generate cache key for pool
func GenSeatPoolKey(supplierName string, slID int64, publisherID int64, publisherDomain string, crType int64) string {
	return fmt.Sprintf(
		"sl%d_sup%s_pubid%d_pubdo%s_crtype%d",
		slID,
		supplierName,
		publisherID,
		publisherDomain,
		crType,
	)
}

// Encode is the encode function for serialize object in io writer
func (s Seat) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(s)
}

// Decode try to decode object from io reader
func (s Seat) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(s)
}
