package manager

import (
	"fmt"
	"strings"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/services/broker"
)

// TableModel is the model for counting data and aggregate them into on query
type TableModel struct {
	Supplier      string // all
	Source        string // all
	Demand        string // all except impression
	Time          int64  // all
	Request       int64  // impression
	Impression    int64  // impression slot
	Win           int64  // demand
	Show          int64  // show
	ImpressionBid int64  // winner
	ShowBid       int64  // show
	WinnerBid     int64  // Winner
	Acknowledger  *broker.Delivery
}

func (tm *TableModel) hasDataSupplierDemand() (bool, string, []interface{}) {
	b := tm.Request+tm.Impression+tm.Win+tm.Show+tm.ImpressionBid+tm.WinnerBid+tm.ShowBid > 0
	if b {
		//(supplier,demand,source,time_id,request_count,impression_count,win_count,show_count,imp_bid,show_bid,win_bid)
		return true, "(?,?,?,?,?,?,?,?,?,?,?)", []interface{}{
			tm.Supplier,
			tm.Demand,
			tm.Source,
			tm.Time,
			tm.Request,
			tm.Impression,
			tm.Win,
			tm.Show,
			tm.ImpressionBid,
			tm.ShowBid,
			tm.WinnerBid,
		}
	}
	return false, "", nil
}

func (tm *TableModel) hasDataSupplier() (bool, string, []interface{}) {
	b := tm.Request+tm.Impression+tm.Show+tm.ImpressionBid+tm.ShowBid > 0
	if b {
		// (supplier,source,time_id,request_count,impression_count,show_count,imp_bid,show_bid)
		return true, "(?,?,?,?,?,?,?,?)", []interface{}{
			tm.Supplier,
			tm.Source,
			tm.Time,
			tm.Request,
			tm.Impression,
			tm.Show,
			tm.ImpressionBid,
			tm.ShowBid,
		}
	}

	return false, "", nil
}

const supDemSrcTable = `INSERT INTO sup_dem_src
(supplier,demand,source,time_id,request_count,impression_count,win_count,show_count,imp_bid,show_bid,win_bid) VALUES
%s
ON DUPLICATE KEY UPDATE
 request=request+VALUES(request), impression=impression+VALUES(impression),win=win+VALUES(win),
 show_time=show_time+VALUES(show_time),imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid),win_bid=win_bid+VALUES(win_bid)
`

const supSrcTable = `INSERT INTO sup_dem_src
(supplier,source,time_id,request_count,impression_count,show_count,imp_bid,show_bid) VALUES
%s
ON DUPLICATE KEY UPDATE
 request_count=request_count+VALUES(request_count), impression_count=impression_count+VALUES(impression_count),
 show_count=show_count+VALUES(show_count),imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid)
`

func flush(supDemSrc map[string]*TableModel, supSrc map[string]*TableModel) error {
	var (
		parts1, parts2   []string
		params1, params2 []interface{}
	)
	for i := range supDemSrc {
		if has, part, param := supDemSrc[i].hasDataSupplierDemand(); has {
			parts1 = append(parts1, part)
			params1 = append(params1, param...)
		}
	}
	q1 := fmt.Sprintf(supDemSrcTable, strings.Join(parts1, ","))

	for i := range supSrc {
		if has, part, param := supDemSrc[i].hasDataSupplier(); has {
			parts2 = append(parts2, part)
			params2 = append(params2, param...)
		}
	}
	q2 := fmt.Sprintf(supSrcTable, strings.Join(parts2, ","))

	return models.NewManager().MultiQuery(
		models.Parts{
			Query:  q1,
			Params: params1,
			Do:     len(params1) > 0,
		},
		models.Parts{
			Query:  q2,
			Params: params2,
			Do:     len(params2) > 0,
		},
	)

}
