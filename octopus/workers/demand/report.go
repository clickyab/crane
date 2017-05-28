package demand

import (
	"fmt"
	"time"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/services/assert"
)

var epoch time.Time

// UpdateReport will update demand report (inclusive)
func UpdateReport(t time.Time) {

	td := t.Format("2006-01-02")
	from, to := factTableRange(t)

	var q = fmt.Sprintf(`INSERT INTO demand_report (
								demand,
								target_date,
								request_out_count,
								imp_in_count,
								imp_out_count,
								win_count,
								win_bid,
								deliver_count,
								deliver_bid,
								profit
								)

							SELECT demand,
							"%s",
							sum(request_out_count),
							sum(imp_in_count),
							sum(imp_out_count),
							sum(win_count),
							sum(win_bid),
							sum(deliver_count),
							sum(deliver_bid),
							sum(profit)
								FROM sup_dem_src WHERE time_id BETWEEN %d AND %d
							GROUP BY demand

							 ON DUPLICATE KEY UPDATE
							  demand=VALUES(demand),
							  target_date=VALUES(target_date),
							  request_out_count=VALUES(request_out_count),
							  imp_in_count=VALUES(imp_in_count),
							  imp_out_count=VALUES(imp_out_count),
							  win_count=VALUES(win_count),
							  win_bid=VALUES(win_bid),
							  deliver_count=VALUES(deliver_count),
							  deliver_bid=VALUES(deliver_bid),
							  profit=values(profit)`, td, from, to)

	_, err := models.NewManager().GetRDbMap().Exec(q)
	assert.Nil(err)
}

// UpdateRange will update demand report in range of two date (inclusive)
func UpdateRange(from time.Time, to time.Time) {
	if from.Unix() > to.Unix() {
		from, to = to, from
	}
	to = to.Add(24 * time.Hour)
	for from.Unix() < to.Unix() {
		UpdateReport(from)
		from = from.Add(time.Hour * 24)
	}
}

// FactTableID is a helper function to get the fact table id from time
func factTableID(tm time.Time) int64 {
	return int64(tm.Sub(epoch).Hours()) + 1
}

func factTableRange(t time.Time) (int64, int64) {
	y, m, d := t.Date()
	from := time.Date(y, m, d, 0, 0, 1, 0, time.UTC)
	to := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	return factTableID(from), factTableID(to)
}

// Report is base model for demand reports
type Report struct {
	ID                  int64  `json:"id"`
	Demand              string `json:"demand"`
	ImpressionOut       int64  `json:"impression_out"`
	SuccessImpression   int64  `json:"success_impression"`
	SuccessRate         int64  `json:"success_rate"`
	WinImpression       int64  `json:"win_impression"`
	WinRate             int64  `json:"win_rate"`
	DeliveredImpression int64  `json:"delivered_impression"`
	DeliverRate         int64  `json:"deliver_rate"`
	Spent               int64  `json:"spent"`
}

func calculator(a []models.DemandReport) []Report {
	res := make([]Report, 0)

	for _, v := range a {
		res = append(res, Report{
			Demand:              v.Demand,
			ID:                  v.ID,
			ImpressionOut:       v.ImpOutCount,
			SuccessImpression:   v.ImpInCount,
			SuccessRate:         (v.ImpOutCount * 100) / v.ImpInCount,
			DeliveredImpression: v.DeliverCount,
			DeliverRate:         (v.DeliverCount * 100) / v.WinCount,
			WinImpression:       v.WinCount,
			WinRate:             (v.WinCount * 100) / v.ImpInCount,
			Spent:               v.DeliverBid,
		})
	}

	return res
}

// ByDate returns list of demand for specific date
func ByDate(t time.Time) []Report {
	return ByDateRange(t, t)
}

// ByDateRange returns list of demand for range of dates
func ByDateRange(from time.Time, to time.Time) []Report {
	return ByDateRangeNames(from, to)
}

// ByDateNames returns demand with specific date
func ByDateNames(f time.Time, demands ...string) []Report {
	return ByDateRangeNames(f, f, demands...)
}

// ByDateRangeNames returns demands with for range of dates
func ByDateRangeNames(f time.Time, t time.Time, names ...string) []Report {

	var a []models.DemandReport

	q := fmt.Sprintf(`SELECT
					id,
					demand,
					target_date,
					request_out_count,
					imp_in_count,
					imp_out_count,
					win_count,
					win_bid,
					deliver_count,
					deliver_bid
				FROM demand_report where %s %s ORDER BY id DESC	`,
		timePartial(true, f, t), demandPartial(false, names...))

	_, err := models.NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

// AggregateByDate returns list of demand for specific date
func AggregateByDate(t time.Time) []Report {
	return AggregateDemandsByDateRange(t, t)
}

// AggregateByDateRange return list of demand for range of dates
func AggregateByDateRange(f time.Time, t time.Time) []Report {
	return AggregateDemandsByDateRange(f, t)

}

// AggregateDemandsByDate return demand with specific date
func AggregateDemandsByDate(f time.Time, demands ...string) []Report {
	return AggregateDemandsByDateRange(f, f, demands...)
}

// AggregateDemandsByDateRange return demands with for range of dates
func AggregateDemandsByDateRange(f time.Time, t time.Time, demands ...string) []Report {

	var a []models.DemandReport

	q := fmt.Sprintf(`SELECT
					demand,
					SUM(target_date) as target_date,
					SUM(request_out_count) as request_out_count ,
					SUM(imp_in_count) as imp_in_count,
					SUM(imp_out_count) as imp_out_count,
					SUM(win_count) as win_count,
					SUM(win_bid) as win_bid,
					SUM(deliver_count) as deliver_count,
					SUM(deliver_bid) as deliver_bid
				FROM demand_report where %s %s GROUP BY demand`,
		timePartial(true, f, t), demandPartial(false, demands...))

	_, err := models.NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

// AggregateAllByDate return all with for range of dates
func AggregateAllByDate(t time.Time) []Report {
	return AggregateAllByDateRange(t, t)
}

// AggregateAllByDateRange return demands with for range of dates
func AggregateAllByDateRange(f time.Time, t time.Time) []Report {

	var a []models.DemandReport

	q := fmt.Sprintf(`SELECT
					"All",
					SUM(target_date) as target_date,
					SUM(request_out_count) as request_out_count ,
					SUM(imp_in_count) as imp_in_count,
					SUM(imp_out_count) as imp_out_count,
					SUM(win_count) as win_count,
					SUM(win_bid) as win_bid,
					SUM(deliver_count) as deliver_count,
					SUM(deliver_bid) as deliver_bid
				FROM demand_report where %s`,
		timePartial(true, f, t))

	_, err := models.NewManager().GetRDbMap().Select(&a, q)
	assert.Nil(err)

	return calculator(a)
}

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	var err error
	epoch, err = time.Parse(layout, str)
	assert.Nil(err)
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
