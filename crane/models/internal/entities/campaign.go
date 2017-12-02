package entities

import (
	"clickyab.com/crane/crane/entity"
)

// Campaign implement entity advertise interface
type Campaign struct {
	ad
}

func (c *Campaign) AppBrands() []string {
	return c.CampaignAppBrand.Array()

}

func (c *Campaign) AppCarriers() []string {
	return c.CampaignAppBrand.Array()
}

func (c *Campaign) WebMobile() bool {
	return c.CampaignWebMobile == 1
}

func (c *Campaign) Web() bool {
	return c.CampaignWeb == 1
}

func (c *Campaign) Hoods() []string {
	return c.CampaignHoods.Array()
}

func (c *Campaign) Isp() []string {
	c.CampaignISP.Array()
}

func (c *Campaign) NetProvider() []string {
	return c.CampaignNetProvider.Array()
}

func (c *Campaign) ID() int64 {
	return c.FCampaignID
}

func (c *Campaign) Name() string {
	return c.FCampaignName.String
}

func (c *Campaign) MaxBID() int64 {
	return c.FCampaignMaxBid
}

func (c *Campaign) Frequency() int {
	return c.FCampaignFrequency
}

func (c *Campaign) Target() entity.Target {
	return entity.Target(c.CampaignNetwork)

}

func (c *Campaign) BlackListPublisher() []string {
	if c.CampaignType == 1 {
		return c.CampaignAppFilter.Array()
	}
	return c.CampaignWebsiteFilter.Array()
}

func (c *Campaign) WhiteListPublisher() []string {
	if c.CampaignType == 1 {
		return c.CampaignPlacement.Array()
	}
	return c.CampaignApp.Array()
}

func (c *Campaign) AllowedOS() []string {
	panic("implement me")
}

func (c *Campaign) Country() []string {
	return c.CampaignCountry.Array()
}

func (c *Campaign) Province() []string {
	return c.CampaignRegion.Array()
}

func (c *Campaign) LatLon() (float64, float64, float64) {
	return c.CampaignLatMap.Float64, c.CampaignLongMap.Float64, c.CampaignRadius.Float64
}

func (c *Campaign) Category() []entity.Category {
	cat := make([]entity.Category, 0)
	for _, v := range c.CampaignCat.Array() {
		cat = append(cat, entity.Category(v))
	}
	return cat
}

func (c *Campaign) Attributes() map[string]interface{} {
	panic("implement me")
}
