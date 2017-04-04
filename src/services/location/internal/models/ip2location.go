package models

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"services/cache"
	"time"
)

// IP2Location struct table ip2location
type IP2Location struct {
	IPFrom      int64          `json:"ip_from" db:"ip_from"`
	IPTo        int64          `json:"ip_to" db:"ip_to"`
	CountryCode sql.NullString `json:"country_code" db:"country_code"`
	CountryName sql.NullString `json:"country_name" db:"country_name"`
	RegionName  sql.NullString `json:"region_name" db:"region_name"`
	CityName    sql.NullString `json:"city_name" db:"city_name"`
}

//GetLocation return the location of an ip from location database
func (m *Manager) GetLocation(ip net.IP) (*IP2Location, error) {
	var res IP2Location
	long, err := ip2long(ip)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("IP_%d", long)
	cc := cache.CreateWrapper(key, &res)
	if err := cache.Hit(key, cc); err == nil {
		return &res, err
	}

	query := `SELECT * FROM ip2location_ir WHERE ip_from <= ? AND ip_to >= ? LIMIT 1`
	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		long,
		long,
	)
	err = cache.Do(cc, 72*time.Hour, err)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// ip2long change ip to integer
func ip2long(ip net.IP) (uint32, error) {
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip2 := ip.To4()
	if ip2 == nil {
		return 0, fmt.Errorf("ipv6? the input was %s", ip.String())
	}
	return binary.BigEndian.Uint32(ip2), nil
}
