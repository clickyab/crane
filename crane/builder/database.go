package builder

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"clickyab.com/gad/models"
)

var (
	irmci   = regexp.MustCompile("(?i)(IR)?(-)?(MCI|TCI|43270|Mobile Communications Company of Iran)$")
	irancel = regexp.MustCompile("(?i)(MTN)?(-)?(irancell|mtn|Iran( )?cell Telecommunications Services Company)$")
	rightel = regexp.MustCompile("(?i)(righ( )?tel(@)?|IRN 20)$") // Some case are like "Rightel | rightel"
)

// SetMobileRelatedParam try to set Mobile related parameters
func SetMobileRelatedParam(u *url.URL) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		brand := u.Query().Get("brand")
		network := u.Query().Get("network")

		carrier := strings.Trim(u.Query().Get("carrier"), "# \n\t")
		if irancel.MatchString(carrier) {
			carrier = "Irancell"
		} else if irmci.MatchString(carrier) {
			carrier = "IR-MCI"
		} else if rightel.MatchString(carrier) {
			carrier = "RighTel"
		}
		//extract App stuff
		o.app.GoogleID = u.Query().Get("GoogleAdvertisingId")
		o.app.AndroidID = u.Query().Get("androidid")
		o.app.AndroidDevice = u.Query().Get("deviceid")
		o.app.SDKVersion, _ = strconv.ParseInt(u.Query().Get("clickyabVersion"), 10, 0)

		mcc, _ := strconv.ParseInt(u.Query().Get("mcc"), 10, 0)
		mnc, _ := strconv.ParseInt(u.Query().Get("mnc"), 10, 0)
		cid, _ := strconv.ParseInt(u.Query().Get("cid"), 10, 0)
		lac, _ := strconv.ParseInt(u.Query().Get("lac"), 10, 0)

		var err error
		o.data.CellLocation, err = models.NewManager().GetCellLocation(mcc, mnc, lac, cid, carrier)
		if err != nil {
			return nil, err
		}
		o.data.PhoneData = models.NewManager().GetPhoneData(brand, carrier, network)

		return o, nil
	}
}
