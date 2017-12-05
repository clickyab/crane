package demand

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"clickyab.com/crane/crane/entity"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/random"
)

func ortbWriter(ctx context.Context, c entity.Context, ss []entity.Seat, w io.Writer) http.Header {
	r := openrtb.SeatBid{}
	for _, v := range ss {
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
	j.Encode(openrtb.BidResponse{
		Currency: c.Currency(),
		ID:       <-random.ID,
		SeatBid:  []openrtb.SeatBid{r},
	})
	h := http.Header{}
	h.Set("content-type", "application/json")
	return h
}
