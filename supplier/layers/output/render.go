package output

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"clickyab.com/crane/openrtb"
)

// RenderBanner try to render to suitable json
func RenderBanner(_ context.Context, w http.ResponseWriter, resp *openrtb.BidResponse, extra string) error {
	final := make(map[string]string)
	for i := range resp.GetSeatbid() {
		if len(resp.GetSeatbid()[i].GetBid()) == 0 {
			continue
		}
		slotID := resp.GetSeatbid()[i].GetBid()[0].GetImpid()
		markup := resp.GetSeatbid()[i].GetBid()[0].GetAdm()
		price := resp.GetSeatbid()[i].GetBid()[0].Price
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
