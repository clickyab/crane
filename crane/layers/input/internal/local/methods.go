package local

import (
	"fmt"
	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/publisher"
	"github.com/clickyab/services/ip2location/client"
	"github.com/clickyab/services/random"
	"github.com/mssola/user_agent"
)

// ExtractSlot ExtractSlot
func ExtractSlot(supplier, publisher string, typ publisher.Platforms, width, height int, uniqueID string, attributes map[string]interface{}) entity.Slot {
	return Slot{
		FID:       fmt.Sprintf("%s/%s/%s/%dx%d/%s", supplier, publisher, typ, width, height, uniqueID),
		attribute: attributes,
		FHeight:   height,
		FWidth:    width,
		FTrackID:  <-random.ID,
		// TODO read from config
		slotCTR: 0.1,
	}
}

func OS(ua string) entity.OS {
	userAgent := user_agent.New(ua)
	return entity.OS{
		Valid:  userAgent.OS() != "",
		Name:   userAgent.OS(),
		Mobile: userAgent.Mobile(),
	}

}

func IP(s string) net.IP {
	ip, _, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}
	return net.ParseIP(ip)
}

func FLocation(ip net.IP) Location {
	l := Location{}
	if ip.String() == "" {
		return l
	}
	ip2loc := client.IP2Location(ip.String())

	return Location{
		FCountry: entity.Country{
			Name:  ip2loc.CountryLong,
			ISO:   ip2loc.CountryShort,
			Valid: ip2loc.CountryShort != "",
		},
		FProvince: entity.Province{
			Name:  ip2loc.Region,
			Valid: ip2loc.Region != "",
		},
		FLatLon: entity.LatLon{
			Valid: false,
		},
	}
}
