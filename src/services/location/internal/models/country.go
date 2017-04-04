package models

import (
	"database/sql"
	"services/cache"
	"time"
)

//CountryInfo struct country info
type CountryInfo struct {
	ID        int64          `id:"id" db:"id"`
	Iso       string         `json:"iso" db:"iso"`
	Name      string         `json:"name" db:"name"`
	NiceName  string         `json:"nicename" db:"nicename"`
	Iso3      sql.NullString `json:"iso3" db:"iso3"`
	NumCode   sql.NullString `json:"numcode" db:"numcode"`
	Phonecode sql.NullString `json:"phonecode" db:"phonecode"`
}

//GetCountry get data country from string
func (m *Manager) GetCountry(name string) (*CountryInfo, error) {
	var country CountryInfo

	wr := cache.CreateWrapper("Country_"+name, &country)
	if err := cache.Hit(wr.String(), wr); err == nil {
		return &country, nil
	}

	query := `SELECT * FROM country WHERE iso = ? LIMIT 1`
	err := m.GetRDbMap().SelectOne(
		&country,
		query,
		name,
	)
	err = cache.Do(wr, 72*time.Hour, err)
	if err != nil {
		return nil, err
	}

	return &country, nil
}
