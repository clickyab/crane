package demand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"clickyab.com/crane/crane/entity"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/random"
)

// Render write openrtb bid-response to writer
func Render(_ context.Context, w http.ResponseWriter, ctx entity.Context) error {
	r := openrtb.SeatBid{}
	for _, v := range ctx.Seats() {
		// What if we have no ad for them?
		if v.WinnerAdvertise() == nil {
			continue
		}
		b := openrtb.Bid{
			ID:    v.ReservedHash(),
			ImpID: v.PublicID(),
			AdMarkup: fmt.Sprintf(
				`<iframe src="%s" width="%d" hight="%d" />`,
				v.ShowURL(),
				v.Width(),
				v.Height(),
			),
			AdID:       fmt.Sprint(v.WinnerAdvertise().ID()),
			H:          v.Height(),
			W:          v.Width(),
			Price:      v.CPM() / ctx.Rate(),
			CampaignID: openrtb.StringOrNumber(fmt.Sprint(v.WinnerAdvertise().Campaign().ID())),
		}
		r.Bid = append(r.Bid, b)
	}
	w.Header().Set("content-type", "application/json")
	j := json.NewEncoder(w)
	return j.Encode(openrtb.BidResponse{
		Currency: ctx.Currency(),
		ID:       <-random.ID,
		SeatBid:  []openrtb.SeatBid{r},
	})
}
