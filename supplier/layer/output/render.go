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
func Render(ctx context.Context, w io.Writer, resp *openrtb.BidResponse) error {
	final := make(map[string]string)
	for i := range resp.SeatBid[0].Bid {
		final[resp.SeatBid[0].Bid[i].ImpID] = strings.Replace(resp.SeatBid[0].Bid[i].AdMarkup, "${AUCTION_PRICE}", fmt.Sprintf("%f", resp.SeatBid[0].Bid[i].Price), -1)
	}

	s, err := json.Marshal(final)
	if err != nil {
		return err
	}
	_, err = w.Write(s)
	return err
}
