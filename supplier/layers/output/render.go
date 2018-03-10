package output

import (
	"context"
	"net/http"

	"encoding/json"

	"fmt"
	"strings"

	"github.com/bsm/openrtb"
)

// RenderBanner try to render to suitable json
func RenderBanner(ctx context.Context, w http.ResponseWriter, resp *openrtb.BidResponse, extra string) error {
	final := make(map[string]string)
	for i := range resp.SeatBid {
		if len(resp.SeatBid[i].Bid) == 0 {
			continue
		}
		slotID := resp.SeatBid[i].Bid[0].ImpID
		markup := resp.SeatBid[i].Bid[0].AdMarkup
		price := resp.SeatBid[i].Bid[0].Price
		if slotID == extra {
			slotID = "m"
		}
		final[slotID] = strings.Replace(markup, "${AUCTION_PRICE}", fmt.Sprintf("%f", price), -1)
	}

	s, err := json.Marshal(final)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(s)
	return err
}
