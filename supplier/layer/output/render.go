package output

import (
	"context"
	"io"

	"encoding/json"

	"fmt"
	"strings"

	"github.com/bsm/openrtb"
)

// RenderBanner try to render to suitable json
func RenderBanner(ctx context.Context, w io.Writer, resp *openrtb.BidResponse, extra string) error {
	final := make(map[string]string)
	for i := range resp.SeatBid[0].Bid {
		slotID := resp.SeatBid[0].Bid[i].ImpID
		if slotID == extra {
			slotID = "m"
		}
		final[slotID] = strings.Replace(resp.SeatBid[0].Bid[i].AdMarkup, "${AUCTION_PRICE}", fmt.Sprintf("%f", resp.SeatBid[0].Bid[i].Price), -1)
	}

	s, err := json.Marshal(final)
	if err != nil {
		return err
	}
	_, err = w.Write(s)
	return err
}
