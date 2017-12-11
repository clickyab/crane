package entities

type slot struct {
	SlotID int64 `json:"slot_id"`
}

// InsertSlotAd insert into slot to slots_ads table
func InsertSlotAd(sid, adid int64) (int64, error) {
	m := NewManager().GetRDbMap()
	mo := &slot{}
	err := m.SelectOne(mo, `SELECT slot_id FROM slots_ads where slot_id=? AND ad_id=? `, sid, adid)
	if err == nil {
		return mo.SlotID, nil
	}
	m = NewManager().GetWDbMap()

	query := `INSERT INTO slots_ads (slot_id, ad_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE
sla_id=LAST_INSERT_ID(sla_id)`
	res, err := m.Exec(
		query,
		sid,
		adid,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindSlotID return slot id for given public-id, slot-size
func FindSlotID(pid string, s int) (int64, error) {
	m := NewManager().GetRDbMap()
	mo := &slot{}
	err := m.SelectOne(mo, `SELECT slot_id FROM slots where slot_pubilc_id=? AND slot_size=? `, pid, s)
	if err == nil {
		return mo.SlotID, nil
	}
	m = NewManager().GetWDbMap()

	q := `INSERT INTO slots (slot_pubilc_id, slot_size) VALUES (?, ?) ON DUPLICATE KEY UPDATE
slot_pubilc_id=VALUES(slot_pubilc_id), slot_size=VALUES(slot_size), slot_id=LAST_INSERT_ID(slot_id)`
	res, err := m.Exec(q, pid, s)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
