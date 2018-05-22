package video

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
)

var (
	// Format is start_time/type/duration{/skip, only for linear}
	shortConfig = config.RegisterString("crane.supplier.vast.preset.short", "start/linear/7/3,end/linear/10/3", "short preset.Format is start_time/type/duration{/skip, only for linear}")
	midConfig   = config.RegisterString("crane.supplier.vast.preset.mid", "start/linear/7/3,end/linear/10/3", "mid preset.Format is start_time/type/duration{/skip, only for linear}")
	longConfig  = config.RegisterString("crane.supplier.vast.preset.long", "start/linear/7/3", "long preset.Format is start_time/type/duration{/skip, only for linear}")
	maxDuration = config.RegisterInt("crane.supplier.vast.max_duration", 60, "maximum duration in sec")
	mime        = config.RegisterString("crane.supplier.vast.mimes.default", "video/mp4,image/png,image/gif,image/jpeg,image/jpg", "comma separated list of accepted types")
)

const (
	short = "short"
	mid   = "mid"
	long  = "long"

	linear    = "linear"
	nonLinear = "nonlinear"
)

func getPreset(l string) string {
	switch strings.ToLower(l) {
	case short, mid, long:
		return strings.ToLower(l)
	}

	il, err := strconv.ParseInt(l, 10, 0)
	if err != nil {
		return mid
	}

	if il < 30 {
		return short
	} else if il < 90 {
		return mid
	}
	return long
}

func notZero(in ...int64) int {
	for _, i := range in {
		if i > 0 {
			return int(i)
		}
	}
	return 0
}

func getSingleSeat(c string) output.Seat {
	// start_time/type/duration{/skip, only for linear}

	all := strings.Split(c, "/")
	if len(all) == 4 && all[1] == linear {
		d, _ := strconv.ParseInt(all[2], 10, 0)
		s, _ := strconv.ParseInt(all[3], 10, 0)
		return output.Seat{Start: all[0], Type: linear, Duration: notZero(d, 7), Skip: notZero(s, 3)}
	}
	if len(all) == 3 {
		d, _ := strconv.ParseInt(all[2], 10, 0)
		return output.Seat{Start: all[0], Type: nonLinear, Duration: notZero(d, 7)}
	}
	return output.Seat{Start: "start", Type: linear, Duration: 7, Skip: 3}
}

func getDataFromConfig(c string) map[string]output.Seat {
	var lin, nonlin int
	res := make(map[string]output.Seat)
	d := strings.Split(c, ",")
	for i := range d {
		t := strings.Trim(d[i], "\n\t ")
		s := getSingleSeat(t)
		if s.Type == linear {
			lin++
			s.IDExtra = fmt.Sprintf("%d1", lin)
		} else {
			nonlin++
			s.IDExtra = fmt.Sprintf("%d0", nonlin)
		}

		res[s.Start] = s
	}

	return res
}

func getSlots(l string) map[string]output.Seat {
	p := getPreset(l)
	switch p {
	case short:
		return getDataFromConfig(shortConfig.String())
	case mid:
		return getDataFromConfig(midConfig.String())
	case long:
		return getDataFromConfig(longConfig.String())
	}
	panic("what? " + p)
}

func getMimes(requsted ...string) []string {
	var (
		res   []string
		mimes []string
	)
	for _, m := range strings.Split(mime.String(), ",") {
		if n := strings.Trim(m, "\n\t "); n != "" {
			mimes = append(mimes, n)
		}
	}

	for i := range requsted {
		n := strings.Trim(requsted[i], "\n\t ")
		if n != "" && array.StringInArray(n, mimes...) {
			res = append(res, n)
		}
	}

	if len(res) == 0 {
		return mimes
	}
	return res
}

// the first map is just an array of map, the key is the start value (from the seat) but the returning map is a bit
// tricky, the key is the id, so we can identify each slot in response
func getImps(r *http.Request, pub entity.Publisher, s map[string]output.Seat, requestedMime ...string) ([]openrtb.Impression, map[string]output.Seat) {
	var (
		res   []openrtb.Impression
		times = make(map[string]output.Seat)
		sec   = secure(r)
		mimes = getMimes(requestedMime...)
	)

	baseID := fmt.Sprintf("%d", pub.ID())

	assert.True(len(mimes) > 0)

	// calculate min cpc and insert in impression ext
	impExt := map[string]interface{}{
		"min_cpc": pub.MinCPC(string(entity.RequestTypeVast)),
	}
	iExt, err := json.Marshal(impExt)
	assert.Nil(err)

	for i := range s {
		li := 0
		if s[i].Type == linear {
			li = 1
		}
		times[baseID+s[i].IDExtra] = s[i]
		imp := openrtb.Impression{
			ID:     baseID + s[i].IDExtra,
			Secure: sec,
			Video: &openrtb.Video{
				SkipMin:     s[i].Skip,
				SkipAfter:   s[i].Skip,
				MinDuration: s[i].Skip,
				MaxDuration: maxDuration.Int(),
				Mimes:       mimes,
				Linearity:   li,
				Protocols:   []int{3}, // Only supporting version 3
				Protocol:    3,
				Ext:         iExt,
			},
			BidFloor: float64(pub.FloorCPM()),
		}

		res = append(res, imp)
	}

	return res, times
}

// secure check openrtb protocol (http/https)
func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
