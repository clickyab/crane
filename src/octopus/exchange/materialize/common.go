package materialize

import "octopus/exchange"

func impressionToMap(imp exchange.Impression) map[string]interface{} {
	return map[string]interface{}{
		"track_id":    imp.TrackID(),
		"ip":          imp.IP(),
		"user_agent":  imp.UserAgent(),
		"source":      sourceToMap(imp.Source()),
		"location":    locationToMap(imp.Location()),
		"attributes":  imp.Attributes(),
		"slots":       slotsToMap(imp.Slots()),
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
		"demand":     demandToMap(ad.Demand()),
		"height":     ad.Height(),
		"id":         ad.ID(),
		"landing":    ad.Landing(),
		"max_cpm":    ad.MaxCPM(),
		"rate":       ad.Rates(),
		"track_id":   ad.TrackID(),
		"url":        ad.URL(),
		"width":      ad.Width(),
		"winner_cpm": ad.WinnerCPM(),
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

func slotsToMap(slots []exchange.Slot) []map[string]interface{} {
	resSlots := make([]map[string]interface{}, len(slots))
	for i := range slots {
		resSlots = append(resSlots, map[string]interface{}{
			"height":   slots[i].Height(),
			"track_id": slots[i].TrackID(),
			"width":    slots[i].Width(),
		},
		)
	}
	return resSlots
}

func winnerToMap(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, slotID string) map[string]interface{} {
	return map[string]interface{}{
		"track_id":    imp.TrackID(),
		"demand_name": dmn.Name(),
		"price":       ad.WinnerCPM(),
		"slot_id":     slotID,
	}
}
