package campaign

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// Desktop checker
type Desktop struct {
}

// Check filter network for desktop
func (*Desktop) Check(c entity.Context, in entity.Campaign) error {
	if c.IsMobile() {
		if in.WebMobile() {
			return nil
		}
		return errors.New("DESKTOP_WEBMOBILE")
	}
	if in.Web() {
		return nil
	}
	return errors.New("NO_DESKTOP_WEBMOBILE")

}
