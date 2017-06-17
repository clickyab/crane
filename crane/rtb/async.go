package rtb

import (
	"context"

	"clickyab.com/exchange/crane/entity"
	"github.com/clickyab/services/store"
)

// AsyncCTR is the selection of ads in async mode
func AsyncCTR(ctx context.Context, imp entity.Impression, ads map[int][]entity.Advertise, ch chan map[string]entity.Advertise) {
	s := store.GetSyncStore()
	go selectCTR(ctx, s, imp, ads, ch)
}
