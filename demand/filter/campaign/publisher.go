package campaign

import (
	"errors"
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// WhiteList checker
type WhiteList struct {
}

// Check check if the publisher is in white list of this or not
func (*WhiteList) Check(impression entity.Context, ad entity.Campaign) error {
	if hasString(true, ad.WhiteListPublisher(), fmt.Sprint(impression.Publisher().ID())) {
		return nil
	}
	return errors.New("WHITELIST")

}

// BlackList checker
type BlackList struct {
}

// Check PublisherBlackList checker
func (*BlackList) Check(impression entity.Context, ad entity.Campaign) error {
	if !hasString(false, ad.BlackListPublisher(), fmt.Sprint(impression.Publisher().ID())) {
		return nil
	}
	return errors.New("BLACKLIST")

}
