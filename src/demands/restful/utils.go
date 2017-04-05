package restful

import (
	"entity"
	"fmt"
	"services/statistic"
	"time"
)

const (
	hour  = time.Hour
	day   = 24 * hour
	week  = 7 * day
	month = 31 * day
)

func impressionToMap(imp entity.Impression) map[string]interface{} {
	tmp := map[string]interface{}{
		"impersion_id": imp.MegaIMP(),
		"client_id":    imp.ClientID(),
		"ip":           imp.IP().String(),
		"user_agent":   imp.UserAgent(),
		"publisher":    publisherToMap(imp.Source()),
		"location":     imp.Location(),
		"os":           imp.OS(),
		"slots":        slotToMap(imp.Slots()),
		"category":     imp.Category(),
	}
	return tmp
}

func publisherToMap(pub entity.Publisher) map[string]interface{} {
	tmp := map[string]interface{}{
		"name":           pub.Name(),
		"floor_cpm":      pub.FloorCPM(),
		"soft_floor_cpm": pub.SoftFloorCPM(),
		"type":           pub.Type(),
		"minimum_cpc":    pub.MinCPC(),
		"accepted_type":  pub.AcceptedTypes(),
		"under_floor":    pub.UnderFloor(),
		"supplier":       supplierToMap(pub.Supplier()),
	}
	return tmp
}

func supplierToMap(pub entity.Supplier) map[string]interface{} {
	tmp := map[string]interface{}{
		"name":           pub.Name(),
		"floor_cpm":      pub.FloorCPM(),
		"soft_floor_cpm": pub.SoftFloorCPM(),
		"accepted_type":  pub.AcceptedTypes(),
	}
	return tmp
}

func slotToMap(s []entity.Slot) []map[string]interface{} {
	res := make([]map[string]interface{}, len(s))
	for i := range s {
		res[i] = map[string]interface{}{
			"width":    s[i].Width(),
			"height":   s[i].Height(),
			"state_id": s[i].StateID(),
		}
	}

	return res
}

func getMonthlyPattern() string {
	return time.Now().Format("200601")
}

func getWeeklyPattern() string {
	t := time.Now()
	_, w := t.ISOWeek()

	return time.Now().Format("2006w") + fmt.Sprintf("%02d", w)

}

func getDailyPattern() string {
	return time.Now().Format("20060102")
}

func getHourlyPattern() string {
	return time.Now().Format("2006010203")
}

func getMinutlyPattern() string {
	return time.Now().Format("200601020304")
}

func incCPM(name string, cpm int64) {
	mp := getMonthlyPattern()
	wp := getWeeklyPattern()
	dp := getDailyPattern()
	hp := getHourlyPattern()
	ip := getMinutlyPattern()
	t := statistic.GetStatisticStore(mp+name, month)
	t.IncSubKey("month", cpm)
	t.IncSubKey("month_count", 1)
	t.IncSubKey(dp, cpm)
	t.IncSubKey(dp+"_count", cpm)
	t = statistic.GetStatisticStore(wp+name, week)
	t.IncSubKey("week", cpm)
	t.IncSubKey(dp, cpm)
	t = statistic.GetStatisticStore(dp+name, day)
	t.IncSubKey("day", cpm)
	t.IncSubKey(hp, cpm)
	t = statistic.GetStatisticStore(hp+name, hour)
	t.IncSubKey("hour", cpm)
	t.IncSubKey(ip, cpm)
}

func realVal(all, count int64) int64 {
	if count > 0 {
		return all / count
	}
	return 0
}
func getCPM(name string) (m, w, d, h, i int64) {
	mp := getMonthlyPattern()
	wp := getWeeklyPattern()
	dp := getDailyPattern()
	hp := getHourlyPattern()
	ip := getMinutlyPattern()
	t := statistic.GetStatisticStore(mp+name, month)
	m, _ = t.Touch("month")
	cc, _ := t.Touch("month_count")
	m = realVal(m, cc)

	t = statistic.GetStatisticStore(wp+name, week)
	w, _ = t.Touch("week")
	cc, _ = t.Touch("week_count")
	w = realVal(w, cc)

	t = statistic.GetStatisticStore(dp+name, day)
	d, _ = t.Touch("day")
	cc, _ = t.Touch("day_count")
	d = realVal(d, cc)

	t = statistic.GetStatisticStore(hp+name, hour)
	h, _ = t.Touch("hour")
	cc, _ = t.Touch("hour_count")
	h = realVal(h, cc)

	i, _ = t.Touch(ip)
	cc, _ = t.Touch(ip + "_count")
	i = realVal(i, cc)

	return
}
