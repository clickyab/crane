package jsonbackend

import (
	"octopus/exchange"
)

func impressionToMap(imp exchange.Impression, ads map[string]exchange.Advertise) map[string]interface{} {
	return map[string]interface{}{
		"track_id":    imp.TrackID(),
		"ip":          imp.IP(),
		"user_agent":  imp.UserAgent(),
		"source":      sourceToMap(imp.Source()),
		"location":    locationToMap(imp.Location()),
		"attributes":  imp.Attributes(),
		"slots":       slotsToMap(imp.Slots(), ads),
		"category":    imp.Category(),
		"platform":    imp.Platform(),
		"under_floor": imp.UnderFloor(),
	}
}

func demandToMap(dmn exchange.Demand) map[string]interface{} {
	return map[string]interface{}{
		"name":                 dmn.Name(),
		"call_rate":            dmn.CallRate(),
		"handicap":             dmn.Handicap(),
		"white_list_countries": dmn.WhiteListCountries(),
		"excluded_suppliers":   dmn.WhiteListCountries(),
	}

}

func advertiseToMap(ad exchange.Advertise) map[string]interface{} {
	return map[string]interface{}{
		"demand":        demandToMap(ad.Demand()),
		"height":        ad.Height(),
		"id":            ad.ID(),
		"landing":       ad.Landing(),
		"max_cpm":       ad.MaxCPM(),
		"rate":          ad.Rates(),
		"track_id":      ad.TrackID(),
		"url":           ad.URL(),
		"width":         ad.Width(),
		"winner_cpm":    ad.WinnerCPM(),
		"slot_track_id": ad.SlotTrackID(),
	}
}

func sourceToMap(pub exchange.Publisher) map[string]interface{} {
	return map[string]interface{}{
		"name":           pub.Name(),
		"soft_floor_cpm": pub.SoftFloorCPM(),
		"floor_cpm":      pub.FloorCPM(),
		"attributes":     pub.Attributes(),
		"supplier":       supplierToMap(pub.Supplier()),
	}
}

func supplierToMap(sup exchange.Supplier) map[string]interface{} {
	return map[string]interface{}{
		"floor_cpm":        sup.FloorCPM(),
		"soft_floor_cpm":   sup.SoftFloorCPM(),
		"name":             sup.Name(),
		"share":            sup.Share(),
		"excluded_demands": sup.ExcludedDemands(),
	}
}

func locationToMap(loc exchange.Location) map[string]interface{} {
	return map[string]interface{}{
		"country":  loc.Country(),
		"province": loc.Province(),
		"lat_lon":  loc.LatLon(),
	}
}

func slotsToMap(slots []exchange.Slot, ads map[string]exchange.Advertise) []map[string]interface{} {
	resSlots := make([]map[string]interface{}, 0)
	for i := range slots {
		data := map[string]interface{}{
			"height":   slots[i].Height(),
			"track_id": slots[i].TrackID(),
			"width":    slots[i].Width(),
			"fallback": slots[i].Fallback(),
		}

		if ads != nil {
			if ad, ok := ads[slots[i].TrackID()]; ok {
				data["ad"] = advertiseToMap(ad)
			}
		}
		resSlots = append(resSlots, data)
	}
	return resSlots
}

func winnerToMap(imp exchange.Impression, ad exchange.Advertise, slotID string) map[string]interface{} {
	return map[string]interface{}{
		"track_id":    imp.TrackID(),
		"demand_name": ad.Demand().Name(),
		"price":       ad.WinnerCPM(),
		"slot_id":     slotID,
	}
}

func showToMap(trackID, demand, slotID, adID string, winner int64, supplier string, publisher string) map[string]interface{} {
	return map[string]interface{}{
		"track_id":    trackID,
		"demand_name": demand,
		"price":       winner,
		"slot_id":     slotID,
		"ad_id":       adID,
		"supplier":    supplier,
		"publisher":   publisher,
	}
}
