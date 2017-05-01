package materialize

import "octopus/exchange"

func impressionToMap(imp exchange.Impression) map[string]interface{} {
	return map[string]interface{}{
		"trackID":    imp.TrackID(),
		"IP":         imp.IP(),
		"userAgent":  imp.UserAgent(),
		"source":     sourceToMap(imp.Source()),
		"location":   locationToMap(imp.Location()),
		"attributes": imp.Attributes(),
		"slots":      slotsToMap(imp.Slots()),
		"category":   imp.Category(),
		"platForm":   imp.Platform(),
		"unserFloor": imp.UnderFloor(),
	}
}

func demandToMap(dmn exchange.Demand) map[string]interface{} {
	return map[string]interface{}{
		"name":               dmn.Name(),
		"callRate":           dmn.CallRate(),
		"handicap":           dmn.Handicap(),
		"whiteListCountries": dmn.WhiteListCountries(),
		"excludedSuppliers":  dmn.WhiteListCountries(),
	}

}

func advertiseToMap(ad exchange.Advertise) map[string]interface{} {
	return nil
}

func sourceToMap(pub exchange.Publisher) map[string]interface{} {
	return map[string]interface{}{
		"name":         pub.Name(),
		"softFloorCPM": pub.SoftFloorCPM(),
		"floorCPM":     pub.FloorCPM(),
		"attributes":   pub.Attributes(),
		"supplier":     supplierToMap(pub.Supplier()),
	}
}

func supplierToMap(sup exchange.Supplier) map[string]interface{} {
	return map[string]interface{}{
		"floorCPM":        sup.FloorCPM(),
		"softFloorCPM":    sup.SoftFloorCPM(),
		"name":            sup.Name(),
		"share":           sup.Share(),
		"excludedDemands": sup.ExcludedDemands(),
	}
}

func locationToMap(loc exchange.Location) map[string]interface{} {
	return map[string]interface{}{
		"country":  loc.Country(),
		"province": loc.Province(),
		"latLon":   loc.LatLon(),
	}
}

func slotsToMap(slots []exchange.Slot) []map[string]interface{} {
	resSlots := make([]map[string]interface{}, len(slots))
	for i := range slots {
		resSlots = append(resSlots, map[string]interface{}{
			"height":  slots[i].Height(),
			"trackID": slots[i].TrackID(),
			"width":   slots[i].Width(),
		},
		)
	}
	return resSlots
}
