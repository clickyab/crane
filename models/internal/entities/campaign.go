package entities

import (
	"fmt"
	"strings"

	"clickyab.com/crane/demand/entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/config"
)

var clickyabNetwork = map[string]networkConn{
	"2g":   cellular2G,
	"edge": cellular2G,
	"gprs": cellular2G,
	"3g":   cellular3G,
	"4g":   cellular4G,
}

type networkConn int

const (
	cellular2G = 4
	cellular3G = 5
	cellular4G = 6
)

// Campaign implement entity advertise interface
type Campaign struct {
	ad
	category       []openrtb.ContentCategory
	strategy       entity.Strategy
	connectionType []int
}

// Strategy of campaign. can be cpm, cpc
func (c *Campaign) Strategy() entity.Strategy {
	if c.strategy > 0 {
		return c.strategy
	}
	switch strings.ToLower(c.CampaignBillingType.String) {
	case "cpm":
		c.strategy = entity.StrategyCPM
	case "cpc":
		c.strategy = entity.StrategyCPC
	default:
		c.strategy = entity.StrategyCPC
	}
	return c.strategy
}

// AppBrands return accepted app brands id
func (c *Campaign) AppBrands() []string {
	return c.CampaignAppBrandName.Array()

}

// AppCarriers return accepted app carriers
func (c *Campaign) AppCarriers() []string {
	return c.CampaignAppsCarriersName.Array()
}

// WebMobile is this is accepted for web mobile
func (c *Campaign) WebMobile() bool {
	return c.CampaignWebMobile == 1
}

// Web is this accepted for web
func (c *Campaign) Web() bool {
	return c.CampaignWeb == 1
}

// Hoods for this request
func (c *Campaign) Hoods() []string {
	return c.CampaignHoods.Array()
}

// ISP return the isp
func (c *Campaign) ISP() []string {
	return c.CampaignISP.Array()
}

// ConnectionType return accepted net providers id
func (c *Campaign) ConnectionType() []int {
	if c.connectionType == nil {
		c.connectionType = make([]int, 0)
		for i := range c.CampaignNetProviderName.Array() {
			val, ok := clickyabNetwork[strings.ToLower(c.CampaignNetProviderName.Array()[i])]
			if ok {
				c.connectionType = append(c.connectionType, int(val))
			}
		}
	}
	return c.connectionType
}

// ID campaign id
func (c *Campaign) ID() int32 {
	return c.FCampaignID
}

// Name campaign name
func (c *Campaign) Name() string {
	return c.FCampaignName.String
}

var minFrequency = config.RegisterInt("crane.models.min_frequency", 2, "min frequency for campaign")

// Frequency campaign frequency
func (c *Campaign) Frequency() int32 {
	if c.FCampaignFrequency <= 0 {
		c.FCampaignFrequency = int32(minFrequency.Int())
	}
	return c.FCampaignFrequency
}

// BlackListPublisher campaign black list publishers
func (c *Campaign) BlackListPublisher() []string {
	if entity.Target(c.CampaignNetwork) == entity.TargetApp {
		return c.CampaignAppFilter.Array()
	}
	return c.CampaignWebsiteFilter.Array()
}

// WhiteListPublisher return the campaign white list publishers
func (c *Campaign) WhiteListPublisher() []string {
	if entity.Target(c.CampaignNetwork) == entity.TargetApp {
		return c.CampaignApp.Array()
	}
	return c.CampaignPlacement.Array()
}

// Country allowed countries
func (c *Campaign) Country() []string {
	return c.CampaignCountry.Array()
}

// Province allowed province (Iran only)
func (c *Campaign) Province() []string {
	return c.CampaignGeos.Array()
}

// LatLon allowed lat lon radius
func (c *Campaign) LatLon() (bool, float64, float64, float64) {
	b := c.CampaignLatMap.Valid && c.CampaignLongMap.Valid && c.CampaignRadius.Valid
	return b, c.CampaignLatMap.Float64, c.CampaignLongMap.Float64, c.CampaignRadius.Float64
}

// Category allowed category
func (c *Campaign) Category() []openrtb.ContentCategory {
	if c.category != nil {
		return c.category
	}
	c.category = make([]openrtb.ContentCategory, 0)
	if !c.CampaignBillingType.Valid {
		return c.category
	}
	for _, v := range c.CampaignCat.Array() {

		r, ok := openrtb.ContentCategory_value["IAB"+v]
		if ok {
			c.category = append(c.category, openrtb.ContentCategory(r))
		}
	}
	fmt.Println("cat: ", c.FCampaignID, c.category)
	return c.category
}

// AllowedOS is the allowed os for this campaign
func (c *Campaign) AllowedOS() []string {
	return c.ad.CampaignPlatforms.Array()
}
