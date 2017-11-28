package models

import (
	"strings"

	"clickyab.com/crane/crane/entity"
)

type campaign struct {
	ad
}

func (c *campaign) ID() int64 {
	return c.FCampaignID
}

func (c *campaign) Name() string {
	return c.FCampaignName.String
}

func (c *campaign) MaxBID() int64 {
	return c.FCampaignMaxBid
}

func (c *campaign) Frequency() int {
	return c.FCampaignFrequency
}

func (c *campaign) Target() entity.Target {
	return entity.Target(c.CampaignNetwork)

}

func (c *campaign) BlackListPublisher() []string {
	if c.CampaignType == 1 {
		return strings.Split(string(c.CampaignAppFilter), "#")
	}
	return strings.Split(string(c.CampaignWebsiteFilter), "#")
}

func (c *campaign) WhiteListPublisher() []string {
	if c.CampaignType == 1 {
		return strings.Split(string(c.CampaignPlacement), "#")
	}
	return strings.Split(string(c.CampaignApp), "#")
}

func (c *campaign) AllowedOS() []string {
	panic("implement me")
}

func (c *campaign) Country() []string {
	return strings.Split(string(c.CampaignCountry), "#")
}

func (c *campaign) Province() []string {
	return strings.Split(string(c.CampaignRegion), "#")
}

func (c *campaign) LatLon() (float64, float64, float64) {
	return c.CampaignLatMap.Float64, c.CampaignLongMap.Float64, c.CampaignRadius.Float64
}

func (c *campaign) Category() []entity.Category {
	cat := make([]entity.Category, 0)
	for _, v := range strings.Split(string(c.CampaignCat), "#") {
		cat = append(cat, entity.Category(v))
	}
	return cat
}

func (c *campaign) Attributes() map[string]interface{} {
	panic("implement me")
}
