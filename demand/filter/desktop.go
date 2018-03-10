package filter

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// Desktop checker
type Desktop struct {
}

// Check filter network for desktop
func (*Desktop) Check(c entity.Context, in entity.Creative) error {
	if c.IsMobile() {
		if in.Campaign().WebMobile() {
			return nil
		}
		return errors.New("desktop campaign not webmobile")
	}
	if in.Campaign().Web() {
		return nil
	}
	return errors.New("desktop campaign is not web or webmobile")

}
