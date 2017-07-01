package vast

import (
	"encoding/xml"
	"time"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/rs/vast"
	"github.com/rs/vmap"
)

type (
	adType       string
	breakType    string
	deliveryType string
)

const (
	vastAttr = "vast"
	// its one cuz we always have one ad
	sequence = 1

	imageAdType adType = "image"

	linearType    breakType = "linear"
	nonLinearType breakType = "non-linear"
	//displayType breakType = "display"

	progressiveDeliveryType deliveryType = "progressive"
	bitrate                              = 24
)

var (
	podMultipleAd = config.RegisterBoolean("crane.vast.pod_multiple_ad", false, "it's true when we have mulltiple ad inside a pod")
	redirection   = config.RegisterBoolean("crane.vast.allow_redirection", true, "it's true cuz we redirect the client")
	skipOffset    = config.RegisterDuration("crane.vast.skip_offset", time.Second*0, "skip offset is the cooldown to show skip ad button which we dont use yet so its zero")
)

type slotVastAttribute struct {
	Duration, Offset time.Duration
	BreakType        breakType
}

// Renderer is a function that gets impression and returns marshaled xml vmap
// adBreakID is impressions.TrackID()
// each adID is ad.ID()
func Renderer(impression entity.Impression, cp entity.ClickProvider) []byte {
	var adBreaks []vmap.AdBreak

	for i := range impression.Slots() {
		attribute, ok := impression.Slots()[i].Attribute()[vastAttr].(*slotVastAttribute)
		assert.True(ok)

		advertise := impression.Slots()[i].WinnerAdvertise()
		slot := impression.Slots()[i]

		// check for linear or non-linear
		var il *vast.InLine
		if attribute.BreakType == linearType {
			il = newInlineLinearAd(slot.ShowURL(), cp.ClickURL(slot, impression), imageAdType, advertise.Width(), advertise.Height())
		} else if attribute.BreakType == nonLinearType {
			il = newInlineNonLinearAd(slot.ShowURL(), cp.ClickURL(slot, impression), imageAdType, advertise.Width(), advertise.Height())
		}

		ad := newAd(il)
		vd := newVastData(*ad)

		a, b := podMultipleAd.Bool(), redirection.Bool()
		source := newAdSourceWithVastAdData(string(advertise.ID()), vd, &a, &b)
		adBreak := newAdBreak(attribute.Offset, attribute.Duration, impression.TrackID(), string(attribute.BreakType), source, []vmap.Tracking{})

		adBreaks = append(adBreaks, adBreak)
	}

	vmap := vmap.VMAP{
		Version:  "3.0",
		AdBreaks: adBreaks,
	}

	b, err := xml.Marshal(vmap)
	assert.Nil(err)

	return b
}

// TODO: event handler is needed
// maxbitrate havent implemented
// one impression cuz no multiple vast ad
func newInlineLinearAd(showURL, clickURL string, adtype adType, width, height int) *vast.InLine {
	duration := vast.Duration(skipOffset.Duration())
	return &vast.InLine{
		AdSystem: &vast.AdSystem{
			Version: "0.1.0",
			Name:    "clickyab",
		},
		Impressions: []vast.Impression{{
			URI: "clickyab.com",
		}},
		Creatives: []vast.Creative{{
			Linear: &vast.Linear{
				//its set to zero
				SkipOffset: &vast.Offset{Duration: &duration},
				Duration:   vast.Duration(duration),
				VideoClicks: &vast.VideoClicks{ClickThroughs: []vast.VideoClick{{
					URI: clickURL,
				}}},
				MediaFiles: []vast.MediaFile{{
					Delivery: string(progressiveDeliveryType),
					Type:     string(adtype),
					Width:    width,
					Height:   height,
					Bitrate:  bitrate,
					URI:      showURL,
				}},
			},
			Sequence: sequence,
		}},
	}
}

// TODO: event handler is needed
// maxbitrate havent implemented
// one impression cuz no multiple vast ad
func newInlineNonLinearAd(showURL, clickURL string, adtype adType, width, height int) *vast.InLine {
	return &vast.InLine{
		AdSystem: &vast.AdSystem{
			Version: "0.1.0",
			Name:    "clickyab",
		},
		Impressions: []vast.Impression{{
			URI: showURL,
		}},
		Creatives: []vast.Creative{{
			NonLinearAds: &vast.NonLinearAds{NonLinears: []vast.NonLinear{{
				Width:                 width,
				Height:                height,
				NonLinearClickThrough: clickURL,
			}}},
			Sequence: sequence,
		}},
	}
}

func newAd(inLine *vast.InLine) *vast.Ad {
	return &vast.Ad{
		InLine: inLine,
	}
}

// it gets only one ad cuz we dont have multiple ad vast
func newVastData(ad vast.Ad) *vast.VAST {
	return &vast.VAST{
		Version: "3.0",
		Ads:     []vast.Ad{ad},
	}
}

func newAdSourceWithVastAdData(id string, data *vast.VAST, allowMultipleAd, allowRedirection *bool) *vmap.AdSource {
	return &vmap.AdSource{
		ID:               id,
		VASTAdData:       data,
		AllowMultipleAds: allowMultipleAd,
		FollowRedirects:  allowRedirection,
	}
}

func newAdBreak(timeOffset, repeatAfter time.Duration, breakID string, breakType string, source *vmap.AdSource, events []vmap.Tracking) vmap.AdBreak {
	duration := vast.Duration(timeOffset)
	return vmap.AdBreak{
		AdSource:       source,
		TimeOffset:     vmap.Offset{Duration: &duration},
		BreakID:        breakID,
		BreakType:      breakType,
		RepeatAfter:    vast.Duration(repeatAfter),
		TrackingEvents: events,
	}
}
