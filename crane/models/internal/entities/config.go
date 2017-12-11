package entities

import (
	"database/sql"

	"github.com/clickyab/services/config"
)

var defaultCTR = config.RegisterFloat64("crane.models.default_ctr", 0.1, "default ctr")

func calc(imp, clk sql.NullInt64) float64 {
	if imp.Int64 < 1000 || clk.Int64 == 0 {
		return defaultCTR.Float64()
	}
	return 0
}
