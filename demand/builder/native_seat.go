package builder

import (
	"net/url"

	"clickyab.com/crane/demand/builder/internal/filters"
	"clickyab.com/crane/demand/entity"
	"github.com/bsm/openrtb/native/request"
)

type nativeSeat struct {
	seat

	filters []entity.Filter
}

func assetToFilterFunc(a []request.Asset) []entity.Filter {
	var res []entity.Filter
	for i := range a {
		f := entity.Filter{
			ID:       a[i].ID,
			Required: a[i].Required != 0,
		}
		if a[i].Image != nil {
			f.Type = entity.AssetTypeImage
			f.SubType = int(a[i].Image.TypeID)
			if a[i].Image.Width > 0 {
				f.Extra = append(f.Extra, filters.ExactWidth(a[i].Image.Width))
			}
			if a[i].Image.Height > 0 {
				f.Extra = append(f.Extra, filters.ExactHeight(a[i].Image.Height))
			}

			if a[i].Image.WidthMin > 0 {
				f.Extra = append(f.Extra, filters.MinWidth(a[i].Image.WidthMin))
			}
			if a[i].Image.HeightMin > 0 {
				f.Extra = append(f.Extra, filters.MinHeight(a[i].Image.HeightMin))
			}

			if len(a[i].Image.Mimes) > 0 {
				f.Extra = append(f.Extra, filters.ContainMimeType(a[i].Image.Mimes...))
			}

		} else if a[i].Title != nil {
			f.Type = entity.AssetTypeText
			f.SubType = entity.AssetTypeTextSubTypeTitle
			if a[i].Title.Length > 0 {
				f.Extra = append(f.Extra, filters.MaxLen(a[i].Title.Length))
			}
		} else if a[i].Data != nil {
			f.Type = entity.AssetTypeText
			f.SubType = int(a[i].Data.TypeID)
			if a[i].Data.Length > 0 {
				f.Extra = append(f.Extra, filters.MaxLen(a[i].Data.Length))
			}

		} else if a[i].Video != nil {
			f.Type = entity.AssetTypeVideo
			f.SubType = entity.AssetTypeVideoSubTypeMain
			if len(a[i].Video.Mimes) > 0 {
				f.Extra = append(f.Extra, filters.ContainMimeType(a[i].Video.Mimes...))
			}
			if a[i].Video.MinDuration > 0 {
				f.Extra = append(f.Extra, filters.MinDuration(a[i].Video.MinDuration))
			}
			if a[i].Video.MaxDuration > 0 {
				f.Extra = append(f.Extra, filters.MaxDuration(a[i].Video.MaxDuration))
			}

			if len(a[i].Video.Protocols) > 0 {
				f.Extra = append(f.Extra, filters.VastProtocol(a[i].Video.Protocols...))
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
		map[string]string{"rh": s.ReservedHash(), "size": "20", "type": s.Type().String(), "subtype": s.SubType().String()},
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
	for i := range s.filters {
		if !s.filters[i].Required {
			continue
		}
		// TODO : think about this. it is somehow complicated thing. what if the asset is selected for two (or more) filter
		// TODO : its a good idea to cache result here, maybe with a key for each AssetFilter function
		if len(advertise.Assets(s.filters[i].Type, s.filters[i].SubType, s.filters[i].Extra...)) > 0 {
			continue
		}
		return false
	}
	return true
}

func (s *nativeSeat) Filters() []entity.Filter {
	return s.filters
}
