package rtb

import (
	"context"
	"crane/entity"
	"services/store"
)

// AsyncCTR is the selection of ads in async mode
func AsyncCTR(ctx context.Context, imp entity.Impression, ads map[int][]entity.Advertise, ch chan map[string]entity.Advertise) {
	s := store.GetSyncStore()
	go selectCTR(ctx, s, imp, ads, ch)
}
