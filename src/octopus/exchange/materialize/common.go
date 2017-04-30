package materialize

import "octopus/exchange"

func impressionToMap(imp exchange.Impression) map[string]interface{} {
	return nil
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
