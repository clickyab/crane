package manager

import (
	"fmt"
	"strings"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
)

func hasDataSupplierDemand(tm *datamodels.TableModel) (bool, string) {
	b := tm.Request+tm.Impression+tm.Win+tm.Show+tm.ImpressionBid+tm.WinnerBid+tm.ShowBid > 0
	if b {
		//(supplier,demand,source,time_id,request_count,impression_count,win_count,show_count,imp_bid,show_bid,win_bid)
		return true, fmt.Sprintf(`("%s", "%s", "%s", %d, %d, %d , %d, %d, %d, %d, %d)`,
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
		)
	}
	return false, ""
}

func hasDataSupplier(tm *datamodels.TableModel) (bool, string) {
	b := tm.Request+tm.Impression+tm.Show+tm.ImpressionBid+tm.ShowBid > 0
	if b {
		// (supplier,source,time_id,request_count,impression_count,show_count,imp_bid,show_bid)
		return true, fmt.Sprintf(`("%s","%s",%d,%d,%d,%d,%d,%d)`,
			tm.Supplier,
			tm.Source,
			tm.Time,
			tm.Request,
			tm.Impression,
			tm.Show,
			tm.ImpressionBid,
			tm.ShowBid,
		)
	}

	return false, ""
}

const supDemSrcTable = `INSERT INTO sup_dem_src
(supplier,demand,source,time_id,request_count,impression_count,win_count,show_count,imp_bid,show_bid,win_bid) VALUES
%s
ON DUPLICATE KEY UPDATE
 request_count=request_count+VALUES(request_count),
 impression_count=impression_count+VALUES(impression_count),
 win_count=win_count+VALUES(win_count),
 show_count=show_count+VALUES(show_count),
 imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid),
 win_bid=win_bid+VALUES(win_bid)
`

const supSrcTable = `INSERT INTO sup_dem_src
(supplier,source,time_id,request_count,impression_count,show_count,imp_bid,show_bid) VALUES
%s
ON DUPLICATE KEY UPDATE
 request_count=request_count+VALUES(request_count),
 impression_count=impression_count+VALUES(impression_count),
 show_count=show_count+VALUES(show_count),
 imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid)
`

func flush(supDemSrc map[string]*datamodels.TableModel, supSrc map[string]*datamodels.TableModel) error {
	var (
		parts1, parts2 []string
	)
	for i := range supDemSrc {
		if has, part := hasDataSupplierDemand(supDemSrc[i]); has {
			parts1 = append(parts1, part)
		}
	}
	q1 := fmt.Sprintf(supDemSrcTable, strings.Join(parts1, ",\n"))

	for i := range supSrc {
		if has, part := hasDataSupplier(supSrc[i]); has {
			parts2 = append(parts2, part)
		}
	}
	q2 := fmt.Sprintf(supSrcTable, strings.Join(parts2, ",\n"))

	return models.NewManager().MultiQuery(
		models.Parts{
			Query: q1,
			Do:    len(parts1) > 0,
		},
		models.Parts{
			Query: q2,
			Do:    len(parts2) > 0,
		},
	)

}
