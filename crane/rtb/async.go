package rtb

import (
	"context"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/kv"
)

// AsyncCTR is the selection of ads in async mode
func AsyncCTR(ctx context.Context, imp entity.Impression, ads map[string][]entity.Advertise, ch chan map[string]entity.Advertise) {
	s := kv.GetSyncStore()
	go selectCTR(ctx, s, imp, ads, ch)
}
