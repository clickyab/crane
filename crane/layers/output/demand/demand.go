package demand

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"clickyab.com/crane/crane/entity"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
)

// Rtb write openrtb bid-response to writer
func Rtb(ctx context.Context, c entity.Context, w io.Writer) http.Header {
	r := openrtb.SeatBid{}
	for _, v := range c.Seats() {
		b := openrtb.Bid{
			ID:         v.ReservedHash(),
			ImpID:      v.PublicID(),
			AdMarkup:   fmt.Sprintf(`<iframe src="%s" w="%d" h="%d" />`, v.ShowURL(), v.Width(), v.Height()),
			AdID:       fmt.Sprint(v.WinnerAdvertise().ID()),
			H:          v.Height(),
			W:          v.Width(),
			Price:      v.Bid(),
			CampaignID: openrtb.StringOrNumber(fmt.Sprint(v.WinnerAdvertise().Campaign().ID())),
		}
		r.Bid = append(r.Bid, b)
	}
	j := json.NewEncoder(w)
	assert.Nil(j.Encode(openrtb.BidResponse{
		Currency: c.Currency(),
		ID:       <-random.ID,
		SeatBid:  []openrtb.SeatBid{r},
	}))
	h := http.Header{}
	h.Set("content-type", "application/json")
	return h
}
