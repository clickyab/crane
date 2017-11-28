package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

//CheckISP find isp
func CheckISP(c *builder.context, in models.AdData) bool {
	if c.GetCommon().ISPID == 0 {
		return len(in.CampaignISP) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignISP.Has(true, c.GetCommon().ISPID)
}
