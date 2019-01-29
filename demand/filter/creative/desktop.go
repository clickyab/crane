package creative

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
		return errors.New("DESKTOP_WEBMOBILE")
	}
	if in.Campaign().Web() {
		return nil
	}
	return errors.New("NO_DESKTOP_WEBMOBILE")

}
