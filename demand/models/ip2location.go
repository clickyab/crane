package models

import (
	"net"
	"regexp"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/ip2location"
)

// IP2lData struct
type IP2lData struct {
	CountryShort string `json:"country_short"`
	CountryLong  string `json:"country_long"`
	Region       string `json:"region"`
	City         string `json:"city"`
	ISP          string `json:"isp"`
}

// IP2Location is the local use of this function
func IP2Location(ip string) IP2lData {
	rec := ip2location.GetAll(ip)
	return IP2lData{
		CountryShort: rec.CountryShort,
		CountryLong:  rec.CountryLong,
		Region:       rec.Region,
		City:         rec.City,
		ISP:          rec.Isp,
	}
}

var ispConst = map[int64]*regexp.Regexp{
	1: regexp.MustCompile(`(?i)iran\s?cell`),
	2: regexp.MustCompile(`(?i)Mobile Communication Company of Iran PLC`),
}

var m = map[string]int64{
	"IR": 1,
	"Azarbayjan-e Sharqi":         2,
	"Ostan-e Azarbayjan-e Gharbi": 3,
	"Ardabil":                     4,
	"Esfahan":                     6,
	"Alborz":                      7,
	"Ilam":                        8,
	"Bushehr":                     9,
	"Tehran":                      10,
	"Chahar Mahall va Bakhtiari":  11,
	"Khorasan-e Janubi":           13,
	"Khorasan-e Razavi":           14,
	"Khorasan-e Shemali":          15,
	"Khuzestan":                   16,
	"Zanjan":                      17,
	"Semnan":                      18,
	"Sistan va Baluchestan":       19,
	"Fars":                        21,
	"Qazvin":                      22,
	"Qom":                         23,
	"Kordestan":                   24,
	"Kerman":                      25,
	"Kermanshah":                  26,
	"Kohkiluyeh va Buyer Ahmadi":  27,
	"Golestan":                    29,
	"Gilan":                       30,
	"Lorestan":                    31,
	"Mazandaran":                  32,
	"Markazi":                     33,
	"Hormozgan":                   34,
	"Hamadan":                     35,
	"Yazd":                        36,
	//"Hamadan":37,
}

// GetProvinceISPByIP get province id by ip
func GetProvinceISPByIP(ip net.IP) entity.Location {
	var province int64
	var uISP int64
	rec := IP2Location(ip.String())
	if i, ok := m[rec.Region]; ok {
		province = i
	}
	if rec.ISP != "" {
		//check isp
		for j := range ispConst {
			if ispConst[j].Match([]byte(rec.ISP)) {
				uISP = j
				break
			}
		}
	}
	l := &location{
		country: entity.Country{
			Name:  rec.CountryLong,
			ISO:   rec.CountryShort,
			Valid: rec.CountryShort != "",
		},
		province: entity.Province{
			Name:  rec.Region,
			Valid: rec.Region != "",
			ID:    province,
		},
		isp: entity.ISP{
			Name:  rec.ISP,
			Valid: rec.ISP != "",
			ID:    uISP,
		},
	}

	return l
}

type location struct {
	country  entity.Country
	province entity.Province
	isp      entity.ISP
	latlon   entity.LatLon
}

func (l *location) Country() entity.Country {
	return l.country
}

func (l *location) Province() entity.Province {
	return l.province
}

func (l *location) LatLon() entity.LatLon {
	return l.latlon
}

func (l *location) ISP() entity.ISP {
	return l.isp
}
