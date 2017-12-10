package filter

import (
	"fmt"

	"clickyab.com/crane/crane/entity"
)

// PublisherWhiteList check if the publisher is in white list of this or not
func PublisherWhiteList(impression entity.Context, ad entity.Advertise) bool {
	return hasString(true, ad.Campaign().WhiteListPublisher(), fmt.Sprint(impression.Publisher().ID()))
}

// PublisherBlackList PublisherBlackList
func PublisherBlackList(impression entity.Context, ad entity.Advertise) bool {
	return !hasString(false, ad.Campaign().BlackListPublisher(), fmt.Sprint(impression.Publisher().ID()))
}
