package entities

import (
	"database/sql"
)

func calc(imp, clk sql.NullInt64) float32 {
	if imp.Int64 < 1000 || clk.Int64 == 0 {
		return -1
	}
	return (float32(clk.Int64) / float32(imp.Int64)) * 10.0
}
