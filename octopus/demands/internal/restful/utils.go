package restful

import (
	"fmt"
	"time"

	"clickyab.com/exchange/services/statistic"
)

const (
	hour  = time.Hour
	day   = 24 * hour
	week  = 7 * day
	month = 31 * day
)

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
