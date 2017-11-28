package rtb

import (
	"time"

	"github.com/clickyab/services/config"
)

var (
	develMode = config.RegisterBoolean("devel_mode", true, "development mode?")

	allowUnderFloor = config.RegisterBoolean("clickyab.under_floor", true, "allow under floor?")

	minCPCVast   = config.RegisterInt("clickyab.min_cpc_vast", 2000, "vast min cpc")
	minCPCApp    = config.RegisterInt("clickyab.min_cpc_app", 700, "app min cpc")
	minCPCNative = config.RegisterInt("clickyab.min_cpc_native", 700, "app min cpc")
	minCPCWeb    = config.RegisterInt("clickyab.min_cpc_native", 2500, "web min cpc")

	floorDivNative = config.RegisterInt("clickyab.floor_div.native", 3, "native floor div")
	floorDivApp    = config.RegisterInt("clickyab.floor_div.app", 1, "app floor div")
	floorDivWeb    = config.RegisterInt("clickyab.floor_div.web", 1, "web floor div")
	floorDivDemand = config.RegisterInt("clickyab.floor_div.demand", 1, "demand floor div")
	floorDivVast   = config.RegisterInt("clickyab.floor_div.vast", 3, "vast floor div")

	minImp = config.RegisterInt("clickyab.min_imp", 1000, "minimum imp")
	//minFrequency     = config.RegisterInt("clickyab.min_frequency", 2, "")
	dailyImpExpire = config.RegisterDuration("clickyab.daily_imp_expire", 7*24*time.Hour, "")
	//dailyClickExpire = config.RegisterDuration("clickyab.daily_click_expire", 7*24*time.Hour, "")
	//dailyCapExpire   = config.RegisterDuration("clickyab.daily_cap_expire", 72*time.Hour, "")
	megaImpExpire = config.RegisterDuration("clickyab.mega_imp_expire", 2*time.Hour, "")
	//convDelay        = config.RegisterDuration("clickyab.conv_delay", time.Second*10, "")
	//convRetry        = config.RegisterInt("clickyab.conv_retry", 8, "")
	minCPMFloorWeb = config.RegisterInt("clickyab.min_cpm_floor_web", 1000, "")
	//minCPMFloorApp   = config.RegisterInt("clickyab.min_cpm_floor_app", 100, "")
	fastClick      = config.RegisterInt("clickyab.fast_click", 2, "")
	adCTREffect    = config.RegisterInt("clickyab.ad_ctr_effect", 70, "")
	slotCTREffect  = config.RegisterInt("clickyab.slot_ctr_effect", 30, "")
	nativeMaxCount = config.RegisterInt("clickyab.native.max_count", 12, "")
	defaultCTR     = config.RegisterFloat64("clickyab.default_ctr", 0.1, "default ctr")
	dailyClickDays = config.RegisterInt("clickyab.daily_click_days", 2, "")

	defaultDuration = config.RegisterString("clickyab.vast.default_duration", "00:00:05", "")
	defaultySkipOff = config.RegisterString("clickyab.vast.default_skipoff", "00:00:03", "")

	chanceShowT = config.RegisterInt("clickyab.chanceshowt", 80, "")
)
