package entities

import "time"

// AddConversion insert to conversion table
func AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID int64, acid string) error {
	q := `INSERT INTO clicks_conv (w_id,app_id,wp_id,ca_id,ad_id,cop_id,cp_id,slot_id,imp_id,c_time,c_date,c_action) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := NewManager().GetWDbMap().Exec(q, wID, appID, wpID, caID, adID, copID, cpID, slotID, impID, time.Now().Unix(), time.Now().Format("20060102"), acid)
	return err
}
