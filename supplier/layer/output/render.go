package output

import (
	"context"
	"io"

	"encoding/json"

	"fmt"
	"strings"

	"github.com/bsm/openrtb"
)

// Render try to render to suitable json
func RenderBanner(ctx context.Context, w io.Writer, resp *openrtb.BidResponse, extra string) error {
	final := make(map[string]string)
	for i := range resp.SeatBid[0].Bid {
		slotId := resp.SeatBid[0].Bid[i].ImpID
		if slotId == extra {
			slotId = "m"
		}
		final[slotId] = strings.Replace(resp.SeatBid[0].Bid[i].AdMarkup, "${AUCTION_PRICE}", fmt.Sprintf("%f", resp.SeatBid[0].Bid[i].Price), -1)
	}

	s, err := json.Marshal(final)
	if err != nil {
		return err
	}
	_, err = w.Write(s)
	return err
}
