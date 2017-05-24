package filter

import "clickyab.com/exchange/crane/entity"

// PublisherWhiteList check if the publisher is in white list of this or not
func PublisherWhiteList(impression entity.Impression, ad entity.Advertise) bool {

	blacklist := ad.WhiteListPublisher()
	if len(blacklist) == 0 {
		// the ad has no black list, pass it by
		return true
	}

	elem := impression.Source().ID()
	return hasInt64(blacklist, elem)
}

// PublisherBlackList PublisherBlackList
func PublisherBlackList(impression entity.Impression, ad entity.Advertise) bool {

	blacklist := ad.BlackListPublisher()
	if len(blacklist) == 0 {
		// no black list is defined. pass it.
		return true
	}

	elem := impression.Source().ID()
	return !hasInt64(blacklist, elem)
}
