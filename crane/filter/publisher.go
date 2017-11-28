package filter

import "clickyab.com/crane/crane/entity"

// PublisherWhiteList check if the publisher is in white list of this or not
func PublisherWhiteList(impression entity.Context, ad entity.Advertise) bool {

	blacklist := ad.Campaign().WhiteListPublisher()
	if len(blacklist) == 0 {
		// the ad has no black list, pass it by
		return true
	}

	elem := impression.Publisher().Name()
	return hasString(blacklist, elem)
}

// PublisherBlackList PublisherBlackList
func PublisherBlackList(impression entity.Context, ad entity.Advertise) bool {

	blacklist := ad.Campaign().BlackListPublisher()
	if len(blacklist) == 0 {
		// no black list is defined. pass it.
		return true
	}

	elem := impression.Publisher().Name()
	return !hasString(blacklist, elem)
}
