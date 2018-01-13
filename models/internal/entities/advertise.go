package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"strconv"
	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

type ad struct {
	FID                      int64                  `db:"ad_id"`
	FType                    int                    `db:"ad_type"`
	FCPM                     int64                  `db:"cpm"`
	FCampaignFrequency       int                    `db:"cp_frequency"`
	FCTR                     float64                `db:"ctr"`
	FCaCTR                   sql.NullFloat64        `db:"ca_ctr"`
	FCampaignMaxBid          int64                  `db:"cp_maxbid"`
	FCampaignID              int64                  `db:"cp_id"`
	FCampaignName            sql.NullString         `db:"cp_name"`
	FAdSize                  int                    `db:"ad_size"`
	FUserID                  int64                  `db:"u_id"`
	FAdName                  sql.NullString         `db:"ad_name"`
	FAdURL                   sql.NullString         `db:"ad_url"`
	FAdCode                  sql.NullString         `db:"ad_code"`
	FAdTitle                 sql.NullString         `db:"ad_title"`
	FAdBody                  sql.NullString         `db:"ad_body"`
	FAdImg                   sql.NullString         `db:"ad_img"`
	FAdStatus                int                    `db:"ad_status"`
	FAdRejectReason          sql.NullString         `db:"ad_reject_reason"`
	FAdConversion            int                    `db:"ad_conv"`
	FAdTime                  int                    `db:"ad_time"`
	FAdMainText              sql.NullString         `db:"ad_mainText"`
	FAdDefineText            sql.NullString         `db:"ad_defineText"`
	FAdTextColor             sql.NullString         `db:"ad_textColor"`
	FAdTarget                sql.NullString         `db:"ad_target"`
	FAdAttribute             mysql.GenericJSONField `db:"ad_attribute"`
	FAdHashAttribute         sql.NullString         `db:"ad_hash_attribute"`
	FCreatedAt               sql.NullString         `db:"created_at"`
	FUpdatedAt               sql.NullString         `db:"updated_at"`
	FUserEmail               string                 `db:"u_email"`
	FUserBalance             string                 `db:"u_balance"`
	FIsCrm                   int                    `db:"is_crm"`
	FCpLock                  int                    `db:"cp_lock"`
	FCampaignAdID            int64                  `db:"ca_id"`
	CampaignType             int                    `db:"cp_type"`
	CampaignBillingType      sql.NullString         `db:"cp_billing_type"`
	CampaignNetwork          int                    `db:"cp_network"`
	CampaignPlacement        SharpArray             `db:"cp_placement"`
	CampaignWebsiteFilter    SharpArray             `db:"cp_wfilter"`
	CampaignRetargeting      sql.NullString         `db:"cp_retargeting"`
	CampaignSegmentID        sql.NullInt64          `db:"cp_segment_id"`
	CampaignNetProvider      SharpArray             `db:"cp_net_provider"`
	CampaignNetProviderName  SharpArray             `db:"cp_net_provider_name"`
	CampaignAppBrand         SharpArray             `db:"cp_app_brand"`
	CampaignAppBrandName     SharpArray             `db:"cp_app_brand_name"`
	CampaignAppLang          sql.NullString         `db:"cp_app_lang"`
	CampaignAppMarket        sql.NullInt64          `db:"cp_app_market"`
	CampaignWebMobile        int                    `db:"cp_web_mobile"`
	CampaignWeb              int                    `db:"cp_web"`
	CampaignApplication      int                    `db:"cp_application"`
	CampaignVideo            int                    `db:"cp_video"`
	CampaignAppsCarriers     SharpArray             `db:"cp_apps_carriers"`
	CampaignAppsCarriersName SharpArray             `db:"cp_app_carrier_name"`
	CampaignLongMap          sql.NullFloat64        `db:"cp_longmap"`
	CampaignLatMap           sql.NullFloat64        `db:"cp_latmap"`
	CampaignRadius           sql.NullFloat64        `db:"cp_radius"`
	CampaignOptCTR           int                    `db:"cp_opt_ctr"`
	CampaignOptConv          int                    `db:"cp_opt_conv"`
	CampaignOptBr            int                    `db:"cp_opt_br"`
	CampaignGender           int                    `db:"cp_gender"`
	CampaignAlexa            int                    `db:"cp_alexa"`
	CampaignFatfinger        int                    `db:"cp_fatfinger"`
	CampaignUnder            int                    `db:"cp_under"`
	CampaignGeos             SharpArray             `db:"cp_geos"`
	CampaignISP              SharpArray             `db:"cp_isp"`
	CampaignRegion           SharpArray             `db:"cp_region"`
	CampaignCountry          SharpArray             `db:"cp_country"`
	CampaignHoods            SharpArray             `db:"cp_hoods"`
	CampaignIspBlacklist     SharpArray             `db:"cp_isp_blacklist"`
	CampaignCat              SharpArray             `db:"cp_cat"`
	CampaignLikeApp          SharpArray             `db:"cp_like_app"`
	CampaignApp              SharpArray             `db:"cp_app"`
	CampaignAppFilter        SharpArray             `db:"cp_app_filter"`
	CampaignKeywords         SharpArray             `db:"cp_keywords"`
	CampaignPlatforms        SharpArray             `db:"cp_platforms"`
	CampaignPlatformVersion  SharpArray             `db:"cp_platform_version"`
	CampaignWeeklyBudget     int                    `db:"cp_weekly_budget"`
	CampaignDailyBudget      int                    `db:"cp_daily_budget"`
	CampaignTotalBudget      int                    `db:"cp_total_budget"`
	CampaignWeeklySpend      int                    `db:"cp_weekly_spend"`
	CampaignTotalSpend       int                    `db:"cp_total_spend"`
	CampaignTodaySpend       int                    `db:"cp_today_spend"`
	CampaignClicks           int                    `db:"cp_clicks"`
	CampaignCTR              float64                `db:"cp_ctr"`
	CampaignImps             int                    `db:"cp_imps"`
	CampaignCPM              int                    `db:"cp_cpm"`
	CampaignCPA              int                    `db:"cp_cpa"`
	CampaignCPC              int                    `db:"cp_cpc"`
	CampaignConv             int                    `db:"cp_conv"`
	CampaignConvRate         float64                `db:"cp_conv_rate"`
	CampaignRevenue          int                    `db:"cp_revenue"`
	CampaignRoi              int                    `db:"cp_roi"`
	CampaignStart            int                    `db:"cp_start"`
	CampaignEnd              int                    `db:"cp_end"`
	CampaignStatus           int                    `db:"cp_status"`
	CampaignLastupdate       int                    `db:"cp_lastupdate"`
	CampaignHourStart        int                    `db:"cp_hour_start"`
	CampaignHourEnd          int                    `db:"cp_hour_end"`
	FMimeType                sql.NullString         `db:"ad_mime"`
}

// AdLoader is the loader of ads
func AdLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []ad
	t := time.Now()
	u := t.Unix()                                                        //return date in unixtimestamp
	h, err := strconv.ParseInt(t.Round(time.Minute).Format("15"), 10, 0) //round time in minute scale
	assert.Nil(err)

	query := fmt.Sprintf(`SELECT
		A.ad_id, C.u_id, ad_name, ad_url,ad_code, ad_title, ad_body, ad_img, ad_status,ad_size,
	 ad_reject_reason, CA.ca_ctr , ad_conv, ad_time, ad_type, ad_mainText, ad_defineText,
	 ad_textColor, ad_target, ad_attribute, ad_hash_attribute, A.created_at, A.updated_at,
	 U.u_email, U.u_balance, C.cp_id, cp_type, cp_billing_type, cp_name, cp_network, cp_placement,
	 cp_wfilter, cp_retargeting, cp_frequency, cp_segment_id, cp_app_brand, cp_net_provider,
	 cp_app_lang, cp_app_market, cp_web_mobile, cp_web, cp_application, cp_video, cp_apps_carriers,
	 cp_longmap, cp_latmap, cp_radius, cp_opt_ctr, cp_opt_conv, cp_opt_br, cp_gender, cp_alexa,
	 cp_fatfinger, cp_under, cp_geos, cp_region, cp_country, cp_hoods, cp_isp_blacklist, cp_cat,
	 cp_like_app, cp_app, cp_app_filter, cp_keywords, cp_platforms, cp_platform_version, cp_maxbid,
	 cp_weekly_budget, cp_daily_budget, cp_total_budget, cp_weekly_spend, cp_total_spend,
	 cp_today_spend, cp_clicks, cp_ctr, cp_imps, cp_cpm, cp_cpa, cp_cpc, cp_conv, cp_conv_rate,
	 cp_revenue, cp_roi, cp_start, cp_end, cp_status, cp_lastupdate, cp_hour_start, cp_hour_end,cp_isp,
	 is_crm, cp_lock,cp_app_brand_name,cp_app_carrier_name,cp_net_provider_name,CA.ca_id, A.ad_mime
	 	FROM campaigns AS C
	 	INNER JOIN users AS U ON C.u_id=U.u_id
		INNER JOIN campaigns_ads AS CA ON C.cp_id=CA.cp_id
		INNER JOIN ads AS A ON A.ad_id=CA.ad_id
		WHERE A.ad_status=1
				AND C.cp_status=1
				AND CA.ca_status = 1
				AND (C.cp_start <= %d OR C.cp_start=0)
				AND (C.cp_end >= %d OR C.cp_end=0)
				AND (cp_time_duration IS NULL OR cp_time_duration LIKE "%%#%d#%%")
				AND C.cp_daily_budget > C.cp_today_spend
				AND C.cp_total_budget > C.cp_total_spend
				AND U.u_balance > U.u_today_spend AND
				U.u_balance > 5000`, u, u, h)

	_, err = NewManager().GetRDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}
	ads := make(map[string]kv.Serializable)
	for i := range res {
		if res[i].FCaCTR.Valid {
			res[i].FCTR = res[i].FCaCTR.Float64
		} else {
			res[i].FCTR = defaultCTR.Float64()
		}
		ads[fmt.Sprint(res[i].FID)] = &Advertise{ad: res[i]}
	}
	return ads, nil
}

// GetAd return a single ad
func GetAd(adID int64) (entity.Advertise, error) {
	query := `SELECT
		A.ad_id, C.u_id, ad_name, ad_url,ad_code, ad_title, ad_body, ad_img, ad_status,ad_size,
	 ad_reject_reason, CA.ca_ctr , ad_conv, ad_time, ad_type, ad_mainText, ad_defineText,
	 ad_textColor, ad_target, ad_attribute, ad_hash_attribute, A.created_at, A.updated_at,
	 U.u_email, U.u_balance, C.cp_id, cp_type, cp_billing_type, cp_name, cp_network, cp_placement,
	 cp_wfilter, cp_retargeting, cp_frequency, cp_segment_id, cp_app_brand, cp_net_provider,
	 cp_app_lang, cp_app_market, cp_web_mobile, cp_web, cp_application, cp_video, cp_apps_carriers,
	 cp_longmap, cp_latmap, cp_radius, cp_opt_ctr, cp_opt_conv, cp_opt_br, cp_gender, cp_alexa,
	 cp_fatfinger, cp_under, cp_geos, cp_region, cp_country, cp_hoods, cp_isp_blacklist, cp_cat,
	 cp_like_app, cp_app, cp_app_filter, cp_keywords, cp_platforms, cp_platform_version, cp_maxbid,
	 cp_weekly_budget, cp_daily_budget, cp_total_budget, cp_weekly_spend, cp_total_spend,
	 cp_today_spend, cp_clicks, cp_ctr, cp_imps, cp_cpm, cp_cpa, cp_cpc, cp_conv, cp_conv_rate,
	 cp_revenue, cp_roi, cp_start, cp_end, cp_status, cp_lastupdate,cp_app_brand_name,cp_app_carrier_name,cp_net_provider_name,cp_hour_start, cp_hour_end,cp_isp,
	 is_crm, cp_lock,CA.ca_id
	 	FROM campaigns AS C
	 	INNER JOIN users AS U ON C.u_id=U.u_id
		INNER JOIN campaigns_ads AS CA ON C.cp_id=CA.cp_id
		INNER JOIN ads AS A ON A.ad_id=CA.ad_id
		WHERE A.ad_id=?`

	res := Advertise{}
	err := NewManager().GetRDbMap().SelectOne(
		&res,
		query,
		adID,
	)
	if err != nil {
		return nil, err
	}
	if res.FCaCTR.Valid {
		res.FCTR = res.FCaCTR.Float64
	} else {
		res.FCTR = defaultCTR.Float64()
	}

	return &res, nil
}

// Advertise implement entity advertise interface
type Advertise struct {
	ad
	campaign entity.Campaign
	size     *size
	capping  entity.Capping
}

// CampaignAdID return campaign_ad primary
func (a *Advertise) CampaignAdID() int64 {
	return a.ad.FCampaignAdID
}

// Size return the size of ad
func (a *Advertise) Size() int {
	return a.ad.FAdSize
}

// Encode is the encode function for serialize object in io writer
func (a *Advertise) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(a.ad)
}

// Decode try to decode object from io reader
func (a *Advertise) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(a.ad)
}

// ID the ad id
func (a *Advertise) ID() int64 {
	return a.FID
}

// Type is the ad type
func (a *Advertise) Type() entity.AdType {
	return entity.AdType(a.FType)
}

// Campaign return the ad campaign object
func (a *Advertise) Campaign() entity.Campaign {
	if a.campaign == nil {
		a.campaign = &Campaign{ad: a.ad}
	}
	return a.campaign
}

// AdCTR return the calculated ad ctr
func (a *Advertise) AdCTR() float64 {
	return a.FCTR
}

// Width is the width of ad
func (a *Advertise) Width() int {
	if a.size == nil {
		a.size = sizes[a.FAdSize]
	}
	return a.size.Width
}

// Height is the height of ad
func (a *Advertise) Height() int {
	if a.size == nil {
		a.size = sizes[a.FAdSize]
	}
	return a.size.Height
}

// Capping return the capping object
func (a *Advertise) Capping() entity.Capping {
	return a.capping
}

// SetCapping is the capping setter
func (a *Advertise) SetCapping(c entity.Capping) {
	a.capping = c
}

// Attributes is the attributes required for dynamic ads
func (a *Advertise) Attributes() map[string]interface{} {
	return a.FAdAttribute
}

// Media return the media address
func (a *Advertise) Media() string {
	if a.FAdImg.Valid {
		return a.FAdImg.String
	}
	return ""
}

// Target return the target address
func (a *Advertise) Target() string {
	if a.FAdURL.Valid {
		return a.FAdURL.String
	}
	return ""
}

// MimeType return the media mime type
func (a *Advertise) MimeType() string {
	return a.FMimeType.String
}
