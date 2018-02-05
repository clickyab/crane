package output

import (
	"context"
	"io"

	"encoding/json"

	"fmt"
	"strings"

	"math/rand"
	"time"

	"github.com/bsm/openrtb"
)

// RenderBanner try to render to suitable json
func RenderBanner(ctx context.Context, w io.Writer, resp *openrtb.BidResponse, extra string, mobile int) error {
	final := make(map[string]string)
	var showT bool
	for i := range resp.SeatBid {
		if len(resp.SeatBid[i].Bid) == 0 {
			continue
		}
		slotID := resp.SeatBid[i].Bid[0].ImpID
		markup := resp.SeatBid[i].Bid[0].AdMarkup
		price := resp.SeatBid[0].Bid[i].Price
		if slotID == extra {
			slotID = "m"
		}
		final[slotID] = strings.Replace(markup, "${AUCTION_PRICE}", fmt.Sprintf("%f", price), -1)
		if mobile == 1 && randomRange(1, 40) == 1 && !showT {
			final[slotID] = fmt.Sprintf(`<iframe width="%d" height="%d" frameborder="0"  scrolling="no">
				%s
				%s
				</iframe>`, resp.SeatBid[i].Bid[0].W, resp.SeatBid[i].Bid[0].H, final[slotID], `<iframe src="//t.clickyab.com/" width="1" height="1" frameborder="0"></iframe>`)
			showT = true
		}
	}

	s, err := json.Marshal(final)
	if err != nil {
		return err
	}
	_, err = w.Write(s)
	return err
}

func randomRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
