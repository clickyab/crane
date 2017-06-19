package console

import (
	"clickyab.com/exchange/octopus/console/internal/routes"
	"github.com/clickyab/services/framework/router"
)

func init() {
	router.Register(&routes.Controller{})
}
