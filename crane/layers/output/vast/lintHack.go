package vast

import (
	"fmt"
	"time"
)

var (
	a = adType("")
	b = breakType("")
	c = deliveryType("")
	d = vastAttr
	e = sequence
	f = imageAdType
	g = linearType
	h = nonLinearType
	i = progressiveDeliveryType
	j = bitrate
)

func init() {
	if false {
		aa := newInlineNonLinearAd("", "", a, 1, 1)
		newInlineLinearAd("", "", a, 1, 1)
		bb := newAd(aa)
		cc := newVastData(*bb)
		aaa, bbb := true, true
		dd := newAdSourceWithVastAdData("", cc, &aaa, &bbb)
		newAdBreak(time.Hour, time.Hour, "", "", dd, nil)
		fmt.Println(a, b, c, d, e, f, g, h, i, j, podMultipleAd, redirection, skipOffset, slotVastAttribute{}, data{})
	}
}
