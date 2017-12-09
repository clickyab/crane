package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"strconv"
	"time"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

type ad struct {
	FID                     int64                  `json:"-" db:"ad_id"`
	FType                   int                    `json:"-" db:"ad_type"`
	FCPM                    int64                  `json:"-" db:"cpm"`
	FCampaignFrequency      int                    `json:"-" db:"cp_frequency"`
	FCTR                    float64                `json:"-" db:"ctr"`
	FCaCTR                  sql.NullFloat64        `json:"-" db:"ca_ctr"`
	FCampaignMaxBid         int64                  `json:"-" db:"cp_maxbid"`
	FCampaignID             int64                  `json:"-" db:"cp_id"`
	FCampaignName           sql.NullString         `json:"-" db:"cp_name"`
	FAdSize                 int                    `json:"-" db:"ad_size"`
	FUserID                 int64                  `json:"-" db:"u_id"`
	FAdName                 sql.NullString         `json:"-" db:"ad_name"`
	FAdURL                  sql.NullString         `json:"-" db:"ad_url"`
	FAdCode                 sql.NullString         `json:"-" db:"ad_code"`
	FAdTitle                sql.NullString         `json:"-" db:"ad_title"`
	FAdBody                 sql.NullString         `json:"-" db:"ad_body"`
	FAdImg                  sql.NullString         `json:"-" db:"ad_img"`
	FAdStatus               int                    `json:"-" db:"ad_status"`
	FAdRejectReason         sql.NullString         `json:"-" db:"ad_reject_reason"`
	FAdConversion           int                    `json:"-" db:"ad_conv"`
	FAdTime                 int                    `json:"-" db:"ad_time"`
	FAdMainText             sql.NullString         `json:"-" db:"ad_mainText"`
	FAdDefineText           sql.NullString         `json:"-" db:"ad_defineText"`
	FAdTextColor            sql.NullString         `json:"-" db:"ad_textColor"`
	FAdTarget               sql.NullString         `json:"-" db:"ad_target"`
	FAdAttribute            mysql.GenericJSONField `json:"-" db:"ad_attribute"`
	FAdHashAttribute        sql.NullString         `json:"-" db:"ad_hash_attribute"`
	FCreatedAt              sql.NullString         `json:"-" db:"created_at"`
	FUpdatedAt              sql.NullString         `json:"-" db:"updated_at"`
	FUserEmail              string                 `json:"-" db:"u_email"`
	FUserBalance            string                 `json:"-" db:"u_balance"`
	FIsCrm                  int                    `json:"-" db:"is_crm"`
	FCpLock                 int                    `json:"-" db:"cp_lock"`
	FCampaignAdID           int64                  `json:"-" db:"ca_id"`
	CampaignType            int                    `json:"-" db:"cp_type"`
	CampaignBillingType     sql.NullString         `json:"-" db:"cp_billing_type"`
	CampaignNetwork         int                    `json:"-" db:"cp_network"`
	CampaignPlacement       SharpArray             `json:"-" db:"cp_placement"`
	CampaignWebsiteFilter   SharpArray             `json:"-" db:"cp_wfilter"`
	CampaignRetargeting     sql.NullString         `json:"-" db:"cp_retargeting"`
	CampaignSegmentID       sql.NullInt64          `json:"-" db:"cp_segment_id"`
	CampaignNetProvider     SharpArray             `json:"-" db:"cp_net_provider"`
	CampaignAppBrand        SharpArray             `json:"-" db:"cp_app_brand"`
	CampaignAppLang         sql.NullString         `json:"-" db:"cp_app_lang"`
	CampaignAppMarket       sql.NullInt64          `json:"-" db:"cp_app_market"`
	CampaignWebMobile       int                    `json:"-" db:"cp_web_mobile"`
	CampaignWeb             int                    `json:"-" db:"cp_web"`
	CampaignApplication     int                    `json:"-" db:"cp_application"`
	CampaignVideo           int                    `json:"-" db:"cp_video"`
	CampaignAppsCarriers    SharpArray             `json:"-" db:"cp_apps_carriers"`
	CampaignLongMap         sql.NullFloat64        `json:"-" db:"cp_longmap"`
	CampaignLatMap          sql.NullFloat64        `json:"-" db:"cp_latmap"`
	CampaignRadius          sql.NullFloat64        `json:"-" db:"cp_radius"`
	CampaignOptCTR          int                    `json:"-" db:"cp_opt_ctr"`
	CampaignOptConv         int                    `json:"-" db:"cp_opt_conv"`
	CampaignOptBr           int                    `json:"-" db:"cp_opt_br"`
	CampaignGender          int                    `json:"-" db:"cp_gender"`
	CampaignAlexa           int                    `json:"-" db:"cp_alexa"`
	CampaignFatfinger       int                    `json:"-" db:"cp_fatfinger"`
	CampaignUnder           int                    `json:"-" db:"cp_under"`
	CampaignGeos            SharpArray             `json:"-" db:"cp_geos"`
	CampaignISP             SharpArray             `json:"-" db:"cp_isp"`
	CampaignRegion          SharpArray             `json:"-" db:"cp_region"`
	CampaignCountry         SharpArray             `json:"-" db:"cp_country"`
	CampaignHoods           SharpArray             `json:"-" db:"cp_hoods"`
	CampaignIspBlacklist    SharpArray             `json:"-" db:"cp_isp_blacklist"`
	CampaignCat             SharpArray             `json:"-" db:"cp_cat"`
	CampaignLikeApp         SharpArray             `json:"-" db:"cp_like_app"`
	CampaignApp             SharpArray             `json:"-" db:"cp_app"`
	CampaignAppFilter       SharpArray             `json:"-" db:"cp_app_filter"`
	CampaignKeywords        SharpArray             `json:"-" db:"cp_keywords"`
	CampaignPlatforms       SharpArray             `json:"-" db:"cp_platforms"`
	CampaignPlatformVersion SharpArray             `json:"-" db:"cp_platform_version"`
	CampaignWeeklyBudget    int                    `json:"-" db:"cp_weekly_budget"`
	CampaignDailyBudget     int                    `json:"-" db:"cp_daily_budget"`
	CampaignTotalBudget     int                    `json:"-" db:"cp_total_budget"`
	CampaignWeeklySpend     int                    `json:"-" db:"cp_weekly_spend"`
	CampaignTotalSpend      int                    `json:"-" db:"cp_total_spend"`
	CampaignTodaySpend      int                    `json:"-" db:"cp_today_spend"`
	CampaignClicks          int                    `json:"-" db:"cp_clicks"`
	CampaignCTR             float64                `json:"-" db:"cp_ctr"`
	CampaignImps            int                    `json:"-" db:"cp_imps"`
	CampaignCPM             int                    `json:"-" db:"cp_cpm"`
	CampaignCPA             int                    `json:"-" db:"cp_cpa"`
	CampaignCPC             int                    `json:"-" db:"cp_cpc"`
	CampaignConv            int                    `json:"-" db:"cp_conv"`
	CampaignConvRate        float64                `json:"-" db:"cp_conv_rate"`
	CampaignRevenue         int                    `json:"-" db:"cp_revenue"`
	CampaignRoi             int                    `json:"-" db:"cp_roi"`
	CampaignStart           int                    `json:"-" db:"cp_start"`
	CampaignEnd             int                    `json:"-" db:"cp_end"`
	CampaignStatus          int                    `json:"-" db:"cp_status"`
	CampaignLastupdate      int                    `json:"-" db:"cp_lastupdate"`
	CampaignHourStart       int                    `json:"-" db:"cp_hour_start"`
	CampaignHourEnd         int                    `json:"-" db:"cp_hour_end"`
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
	 is_crm, cp_lock,CA.ca_id
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
	 cp_revenue, cp_roi, cp_start, cp_end, cp_status, cp_lastupdate, cp_hour_start, cp_hour_end,cp_isp,
	 is_crm, cp_lock,CA.ca_id
	 	FROM campaigns AS C
	 	INNER JOIN users AS U ON C.u_id=U.u_id
		INNER JOIN campaigns_ads AS CA ON C.cp_id=CA.cp_id
		INNER JOIN ads AS A ON A.ad_id=CA.ad_id
		WHERE A.ad_id=$1`

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
