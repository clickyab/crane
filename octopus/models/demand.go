package models

import (
	"fmt"
	"time"

	"github.com/clickyab/services/assert"
)

// updateDemandReport will update demand report (inclusive)
func updateDemandReport(t time.Time) {

	td := t.Format("2006-01-02")
	from, to := factTableRange(t)

	var q = fmt.Sprintf(`INSERT INTO demand_report (
								demand,
								target_date,
								request_out_count,
								ad_in_count,
								imp_out_count,
								ad_out_count,
								ad_out_bid,
								deliver_count,
								deliver_bid,
								profit
								)

							SELECT demand,
							"%s",
							sum(request_out_count),
							sum(ad_in_count),
							sum(imp_out_count),
							sum(ad_out_count),
							sum(ad_out_bid),
							sum(deliver_count),
							sum(deliver_bid),
							sum(profit)
								FROM sup_dem_src WHERE time_id BETWEEN %d AND %d
							GROUP BY demand

							 ON DUPLICATE KEY UPDATE
							  demand=VALUES(demand),
							  target_date=VALUES(target_date),
							  request_out_count=VALUES(request_out_count),
							  ad_in_count=VALUES(ad_in_count),
							  imp_out_count=VALUES(imp_out_count),
							  ad_out_count=VALUES(ad_out_count),
							  ad_out_bid=VALUES(ad_out_bid),
							  deliver_count=VALUES(deliver_count),
							  deliver_bid=VALUES(deliver_bid),
							  profit=values(profit)`, td, from, to)

	_, err := NewManager().GetRDbMap().Exec(q)
	assert.Nil(err)
}

// UpdateDemandRange will update demand report in range of two date (inclusive)
func UpdateDemandRange(from time.Time, to time.Time) {
	if from.Unix() > to.Unix() {
		from, to = to, from
	}
	to = to.Add(24 * time.Hour)
	for from.Unix() < to.Unix() {
		updateDemandReport(from)
		from = from.Add(time.Hour * 24)
	}
}

func calculator(a []DemandReport) []DemandReport {
	res := make([]DemandReport, 0)

	for _, v := range a {
		res = append(res, DemandReport{
			Demand:          v.Demand,
			ID:              v.ID,
			ImpOutCount:     v.ImpOutCount,
			RequestOutCount: v.AdInCount,
			SuccessRate:     (v.ImpOutCount * 100) / v.AdInCount,
			DeliverCount:    v.DeliverCount,
			DeliverRate:     (v.DeliverCount * 100) / v.AdOutCount,
			AdOutCount:      v.AdOutCount,
			WinRate:         (v.AdOutCount * 100) / v.AdInCount,
			DeliverBid:      v.DeliverBid,
		})
	}

	return res
}

// ByDate returns list of demand for specific date
func ByDate(t time.Time) []DemandReport {
	return ByDateRange(t, t)
}

// ByDateRange returns list of demand for range of dates
func ByDateRange(from time.Time, to time.Time) []DemandReport {
	return ByDateRangeNames(from, to)
}

// ByDateNames returns demand with specific date
func ByDateNames(f time.Time, demands ...string) []DemandReport {
	return ByDateRangeNames(f, f, demands...)
}

// ByDateRangeNames returns demands with for range of dates
func ByDateRangeNames(f time.Time, t time.Time, names ...string) []DemandReport {

	var a []DemandReport

	q := fmt.Sprintf(`SELECT
					id,
					demand,
					target_date,
					request_out_count,
					ad_in_count,
					imp_out_count,
					ad_out_count,
					ad_out_bid,
					deliver_count,
					deliver_bid
				FROM demand_report where %s %s ORDER BY id DESC	`,
		timePartial(true, f, t), demandPartial(false, names...))

	_, err := NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

// AggregateByDate returns list of demand for specific date
func AggregateByDate(t time.Time) []DemandReport {
	return AggregateDemandsByDateRange(t, t)
}

// AggregateByDateRange return list of demand for range of dates
func AggregateByDateRange(f time.Time, t time.Time) []DemandReport {
	return AggregateDemandsByDateRange(f, t)

}

// AggregateDemandsByDate return demand with specific date
func AggregateDemandsByDate(f time.Time, demands ...string) []DemandReport {
	return AggregateDemandsByDateRange(f, f, demands...)
}

// AggregateDemandsByDateRange return demands with for range of dates
func AggregateDemandsByDateRange(f time.Time, t time.Time, demands ...string) []DemandReport {

	var a []DemandReport

	q := fmt.Sprintf(`SELECT
					demand,
					target_date,
					SUM(request_out_count) as request_out_count ,
					SUM(ad_in_count) as ad_in_count,
					SUM(imp_out_count) as imp_out_count,
					SUM(ad_out_count) as win_count,
					SUM(ad_out_bid) as win_bid,
					SUM(deliver_count) as deliver_count,
					SUM(deliver_bid) as deliver_bid
				FROM demand_report where %s %s GROUP BY demand`,
		timePartial(true, f, t), demandPartial(false, demands...))

	_, err := NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

// AggregateAllByDate return all with for range of dates
func AggregateAllByDate(t time.Time) []DemandReport {
	return AggregateAllByDateRange(t, t)
}

// AggregateAllByDateRange return demands with for range of dates
func AggregateAllByDateRange(f time.Time, t time.Time) []DemandReport {

	var a []DemandReport

	q := fmt.Sprintf(`SELECT
					"All",
					target_date,
					SUM(request_out_count) as request_out_count ,
					SUM(ad_in_count) as ad_in_count,
					SUM(imp_out_count) as imp_out_count,
					SUM(ad_out_count) as ad_out_count,
					SUM(ad_out_bid) as ad_out_bid,
					SUM(deliver_count) as deliver_count,
					SUM(deliver_bid) as deliver_bid
				FROM demand_report where %s`,
		timePartial(true, f, t))

	_, err := NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

func demandPartial(isFirst bool, names ...string) (res string) {
	if len(names) == 0 {
		return
	}
	if isFirst {
		res = " demand = "
	} else {
		res = "AND demand = "
	}

	for i := range names {
		res += fmt.Sprintf(`"%s"`, names[i])
		if len(names) < i+1 {
			res += " OR "
		}
	}
	return
}

func timePartial(isFirst bool, from time.Time, to time.Time) (res string) {
	if isFirst {
		res = "target_date  "
	} else {
		res = " AND target_date  "
	}
	if from.Unix() > to.Unix() {
		from, to = to, from
	}
	f, e := from.Format("2006-01-02"), to.Format("2006-01-02")
	if f == e {
		res += fmt.Sprintf(` = "%s"`, f)
	} else {
		res += fmt.Sprintf(` BETWEEN "%s" AND "%s"`, f, e)
	}
	return
}
