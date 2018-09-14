package entities

// XXX: Be aware of `INSERT ON DUPLICATE`, just use it after trying normal `SELECT`, Do not use
// Normal Insert

// FindSlotAd insert into slot to slots_ads table
func FindSlotAd(sid int64, adid int32) (int64, error) {
	m := NewManager().GetRDbMap()
	slaID, err := m.SelectInt(`SELECT slot_id FROM slots_ads where slot_id=? AND ad_id=? `, sid, adid)
	if err == nil && slaID != 0 {
		return slaID, nil
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

// FindWebSlotID return slot id for given public-id, slot-size
func FindWebSlotID(pid string, wid int64, s int32) (int64, error) {
	m := NewManager().GetRDbMap()
	// The pubilc is correct. typo is in database
	slID, err := m.SelectInt(`SELECT slot_id FROM slots where slot_pubilc_id=? AND w_id=? `, pid, wid)
	if err == nil && slID != 0 {
		return slID, nil
	}
	m = NewManager().GetWDbMap()

	q := `INSERT INTO slots (slot_pubilc_id,w_id, slot_size) VALUES (?,?, ?) ON DUPLICATE KEY UPDATE
slot_pubilc_id=VALUES(slot_pubilc_id), slot_size=VALUES(slot_size),w_id=VALUES(w_id),
slot_id=LAST_INSERT_ID(slot_id)`
	res, err := m.Exec(q, pid, wid, s)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindAppSlotID return slot id for given public-id, slot-size
func FindAppSlotID(pid string, appid int64, s int32) (int64, error) {
	m := NewManager().GetRDbMap()
	// The pubilc is correct. typo is in database
	slID, err := m.SelectInt(`SELECT slot_id FROM slots where slot_pubilc_id=? AND app_id=? `, pid, appid)
	if err == nil && slID != 0 {
		return slID, nil
	}
	m = NewManager().GetWDbMap()

	q := `INSERT INTO slots (slot_pubilc_id,app_id, slot_size) VALUES (?,?, ?) ON DUPLICATE KEY UPDATE
slot_pubilc_id=VALUES(slot_pubilc_id), slot_size=VALUES(slot_size),app_id=VALUES(app_id),
slot_id=LAST_INSERT_ID(slot_id)`
	res, err := m.Exec(q, pid, appid, s)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
