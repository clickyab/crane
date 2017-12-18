package builder

import "github.com/clickyab/services/config"

var (
	chanceShowT = config.RegisterInt("crane.context.chanceshowt", 80, "chance for showing ad")
)
