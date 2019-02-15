package rtb

import (
	"context"
	"sort"

	"clickyab.com/crane/models/campaign"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
)

// Select is the main entry point for this module
func Select(c context.Context, sel []reducer.Filter, opt ...builder.ShowOptionSetter) (entity.Context, error) {
	// Build context
	ctx, err := builder.NewContext(opt...)
	if err != nil {
		return nil, err
	}
	// Apply filters
	// TODO : after selector fix it
	all, err := reducer.Apply(c, ctx, campaign.GetCampaigns(), sel)
	if err != nil {
		return nil, err
	}

	// select all
	selectAds(c, ctx, all)

	return ctx, nil
}

// MinimalSelect return all ads sorted
func MinimalSelect(
	_ context.Context, ctx *builder.Context, seat entity.Seat, cps []entity.Campaign) (
	[]entity.SelectedCreative, []entity.SelectedCreative) {

	ads := make([]entity.Creative, 0)
	for e := range cps {
		if ak, ok := cps[e].Sizes()[seat.Size()]; ok {
			ads = append(ads, ak...)
		}
	}
	exceed, underfloor := selector(ctx, ads, seat, false, nil)

	exceed = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), exceed)
	underfloor = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), underfloor)

	ef := byMulti{
		Ads:   exceed,
		Video: false,
	}
	uf := byMulti{
		Ads:   underfloor,
		Video: false,
	}
	sort.Sort(ef)
	sort.Sort(uf)

	return ef.Ads, uf.Ads
}
