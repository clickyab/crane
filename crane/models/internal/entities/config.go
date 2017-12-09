package entities

import "github.com/clickyab/services/config"

var defaultCTR = config.RegisterFloat64("crane.models.default_ctr", 0.1, "default ctr")

func calc(imp, clk int) float64 {
	if imp < 1000 || clk == 0 {
		return defaultCTR.Float64()
	}
	return float64(clk) / float64(imp) * 100
}
