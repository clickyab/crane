package entities

import (
	"strings"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/config"
)

// Campaign implement entity advertise interface
type Campaign struct {
	ad
	strategy entity.Strategy
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

// NetProvider return accepted net providers id
func (c *Campaign) NetProvider() []string {
	return c.CampaignNetProviderName.Array()
}

// ID campaign id
func (c *Campaign) ID() int64 {
	return c.FCampaignID
}

// Name campaign name
func (c *Campaign) Name() string {
	return c.FCampaignName.String
}

var minFrequency = config.RegisterInt("crane.models.min_frequency", 2, "min frequency for campaign")

// Frequency campaign frequency
func (c *Campaign) Frequency() int {
	if c.FCampaignFrequency <= 0 {
		c.FCampaignFrequency = minFrequency.Int()
	}
	return c.FCampaignFrequency
}

// BlackListPublisher campaign black list publishers
func (c *Campaign) BlackListPublisher() []string {
	if c.CampaignType == 1 {
		return c.CampaignAppFilter.Array()
	}
	return c.CampaignWebsiteFilter.Array()
}

// WhiteListPublisher return the campaign white list publishers
func (c *Campaign) WhiteListPublisher() []string {
	if c.CampaignType == 1 {
		return c.CampaignPlacement.Array()
	}
	return c.CampaignApp.Array()
}

// Country allowed countries
func (c *Campaign) Country() []string {
	return c.CampaignCountry.Array()
}

// Province allowed province (Iran only)
func (c *Campaign) Province() []string {
	return c.CampaignRegion.Array()
}

// LatLon allowed lat lon radius
func (c *Campaign) LatLon() (bool, float64, float64, float64) {
	b := c.CampaignLatMap.Valid && c.CampaignLongMap.Valid && c.CampaignRadius.Valid
	return b, c.CampaignLatMap.Float64, c.CampaignLongMap.Float64, c.CampaignRadius.Float64
}

// Category allowed category
func (c *Campaign) Category() []entity.Category {
	cat := make([]entity.Category, 0)
	for _, v := range c.CampaignCat.Array() {
		cat = append(cat, entity.Category(v))
	}
	return cat
}

// AllowedOS is the allowed os for this campaign
func (c *Campaign) AllowedOS() []string {
	return c.ad.CampaignPlatforms.Array()
}
