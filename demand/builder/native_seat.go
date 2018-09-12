package builder

import (
	"net/url"

	"clickyab.com/crane/demand/builder/internal/filters"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/openrtb"
)

type nativeSeat struct {
	seat

	filters []entity.Filter
}

func assetToFilterFunc(a []*openrtb.NativeRequest_Asset) []entity.Filter {
	var res []entity.Filter
	for i := range a {
		f := entity.Filter{
			ID:       a[i].GetId(),
			Required: a[i].GetRequired(),
		}
		if a[i].GetImg() != nil {
			f.Type = entity.AssetTypeImage
			f.SubType = int(a[i].GetImg().Type)
			if a[i].GetImg().GetW() > 0 {
				f.Extra = append(f.Extra, filters.ExactWidth(a[i].GetImg().GetW()))
			}
			if a[i].GetImg().GetH() > 0 {
				f.Extra = append(f.Extra, filters.ExactHeight(a[i].GetImg().GetH()))
			}

			if a[i].GetImg().GetWmin() > 0 {
				f.Extra = append(f.Extra, filters.MinWidth(a[i].GetImg().GetWmin()))
			}
			if a[i].GetImg().GetHmin() > 0 {
				f.Extra = append(f.Extra, filters.MinHeight(a[i].GetImg().GetHmin()))
			}

			if len(a[i].GetImg().GetMimes()) > 0 {
				f.Extra = append(f.Extra, filters.ContainMimeType(a[i].GetImg().GetMimes()...))
			}

		} else if a[i].GetTitle() != nil {
			f.Type = entity.AssetTypeText
			f.SubType = entity.AssetTypeTextSubTypeTitle
			if a[i].GetTitle().GetLen() > 0 {
				f.Extra = append(f.Extra, filters.MaxLen(a[i].GetTitle().GetLen()))
			}
		} else if a[i].GetData() != nil {
			f.Type = entity.AssetTypeText
			f.SubType = int(a[i].GetData().GetType())
			if a[i].GetData().GetLen() > 0 {
				f.Extra = append(f.Extra, filters.MaxLen(a[i].GetData().GetLen()))
			}

		} else if a[i].GetVideo() != nil {
			f.Type = entity.AssetTypeVideo
			f.SubType = entity.AssetTypeVideoSubTypeMain
			if len(a[i].GetVideo().Mimes) > 0 {
				f.Extra = append(f.Extra, filters.ContainMimeType(a[i].GetVideo().Mimes...))
			}
			if a[i].GetVideo().GetMinduration() > 0 {
				f.Extra = append(f.Extra, filters.MinDuration(a[i].GetVideo().GetMinduration()))
			}
			if a[i].GetVideo().GetMaxduration() > 0 {
				f.Extra = append(f.Extra, filters.MaxDuration(a[i].GetVideo().GetMaxduration()))
			}

			if len(a[i].GetVideo().Protocols) > 0 {
				f.Extra = append(f.Extra, filters.VastProtocol(a[i].GetVideo().GetProtocols()...))
			}
		}
		res = append(res, f)
	}

	return res
}

// @override
func (s *nativeSeat) ImpressionURL() *url.URL {
	if s.imp != nil {
		return s.imp
	}
	if s.winnerAd == nil {
		panic("no winner")
	}

	s.imp = s.makeURL(
		"pixel",
		map[string]string{
			"rh":      s.ReservedHash(),
			"size":    "20",
			"type":    s.Type().String(),
			"subtype": s.RequestType().String(),
			"pt":      s.context.publisher.Type().String(),
		},
		s.cpm,
		showExpire.Duration(),
	)
	return s.imp
}

// @override
func (s *nativeSeat) Acceptable(advertise entity.Creative) bool {
	if !s.genericTests(advertise) {
		return false
	}

	// check campaign network
	switch s.context.Publisher().Type() {
	case entity.PublisherTypeApp:
		// TODO : fix when implement app native
		return false
	case entity.PublisherTypeWeb:
		// TODO: not totally sure
		if advertise.Target() != entity.TargetNative && advertise.Target() != entity.TargetWeb {
			return false
		}
	default:
		panic("invalid type")
	}

	for i := range s.filters {
		if !s.filters[i].Required {
			continue
		}
		// TODO : think about this. it is somehow complicated thing. what if the asset is selected for two (or more) filter
		// TODO : its a good idea to cache result here, maybe with a key for each AssetFilter function
		if len(advertise.Asset(s.filters[i].Type, s.filters[i].SubType, s.filters[i].Extra...)) > 0 {
			continue
		}
		return false
	}
	return true
}

func (s *nativeSeat) Filters() []entity.Filter {
	return s.filters
}
