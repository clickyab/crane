package entities

import (
	"database/sql"
)

func calc(imp, clk sql.NullInt64) float64 {
	if imp.Int64 < 1000 || clk.Int64 == 0 {
		return -1
	}
	return 0
}
