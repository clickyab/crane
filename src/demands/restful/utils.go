package restful

import "entity"

func impressionToMap(imp entity.Impression) map[string]interface{} {
	tmp := map[string]interface{}{
		"impersion_id": imp.MegaIMP(),
		"client_id":    imp.ClientID(),
		"ip":           imp.IP().String(),
		"user_agent":   imp.UserAgent(),
		"publisher":    publisherToMap(imp.Source()),
		"location":     imp.Location(),
		"os":           imp.OS(),
		"slots":        slotToMap(imp.Slots()),
		"category":     imp.Category(),
	}
	return tmp
}

func publisherToMap(pub entity.Publisher) map[string]interface{} {
	tmp := map[string]interface{}{
		"name":           pub.Name(),
		"floor_cpm":      pub.FloorCPM(),
		"soft_floor_cpm": pub.SoftFloorCPM(),
		"type":           pub.Type(),
		"minimum_cpc":    pub.MinCPC(),
		"accepted_type":  pub.AcceptedTypes(),
		"under_floor":    pub.UnderFloor(),
		"supplier":       supplierToMap(pub.Supplier()),
	}
	return tmp
}

func supplierToMap(pub entity.Supplier) map[string]interface{} {
	tmp := map[string]interface{}{
		"name":           pub.Name(),
		"floor_cpm":      pub.FloorCPM(),
		"soft_floor_cpm": pub.SoftFloorCPM(),
		"accepted_type":  pub.AcceptedTypes(),
	}
	return tmp
}

func slotToMap(s []entity.Slot) map[string]interface{} {
	res := make([]map[string]interface{}, len(s))
	for i := range s {
		res[i] = map[string]interface{}{
			"width":    s[i].Width(),
			"height":   s[i].Height(),
			"state_id": s[i].StateID(),
		}
	}

	return res
}
