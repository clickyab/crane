package filter

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// WhiteList checker
type WhiteList struct {
}

// Check check if the publisher is in white list of this or not
func (*WhiteList) Check(impression entity.Context, ad entity.Creative) bool {
	return hasString(true, ad.Campaign().WhiteListPublisher(), fmt.Sprint(impression.Publisher().ID()))
}

// BlackList checker
type BlackList struct {
}

// Check PublisherBlackList checker
func (*BlackList) Check(impression entity.Context, ad entity.Creative) bool {
	return !hasString(false, ad.Campaign().BlackListPublisher(), fmt.Sprint(impression.Publisher().ID()))
}
