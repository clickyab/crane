package models

import (
	"services/cache"
	"time"
)

//Province struct province info
type Province struct {
	ID     int64  `id:"location_id" db:"location_id"`
	Name   string `json:"location_name" db:"location_name"`
	NameFa string `json:"location_name_persian" db:"location_name_persian"`
	Master int    `json:"location_master" db:"location_master"`
	Select int    `json:"location_select" db:"location_select"`
}

//GetProvince get data province from string
func (m *Manager) GetProvince(name string) (*Province, error) {
	var province Province

	wr := cache.CreateWrapper("Province_"+name, &province)
	if err := cache.Hit(wr.String(), wr); err == nil {
		return &province, nil
	}

	query := `SELECT * FROM list_locations WHERE location_name = ? LIMIT 1`
	err := m.GetRDbMap().SelectOne(
		&province,
		query,
		name,
	)
	err = cache.Do(wr, 72*time.Hour, err)
	if err != nil {
		return nil, err
	}

	return &province, nil
}
