package manager

import (
	"fmt"
	"strings"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
)

func hasDataSupplierDemand(tm *datamodels.TableModel) (bool, string) {
	b := tm.RequestOutCount+tm.ImpressionOutCount+tm.ImpressionInCount+tm.WinCount+tm.WinBid+tm.DeliverCount+tm.DeliverBid+tm.Profit > 0
	if b {
		//(supplier,demand,source,time_id,request_out_count,imp_out_count,imp_in_count,win_count,win_bid,deliver_count,deliver_bid,profit)
		return true, fmt.Sprintf(`("%s", "%s", "%s", %d, %d, %d, %d, %d, %d, %d,%d,%d)`,
			tm.Supplier,
			tm.Demand,
			tm.Source,
			tm.Time,
			tm.RequestOutCount,
			tm.ImpressionOutCount,
			tm.ImpressionInCount,
			tm.WinCount,
			tm.WinBid,
			tm.DeliverCount,
			tm.DeliverBid,
			tm.Profit,
		)
	}
	return false, ""
}

func hasDataSupplier(tm *datamodels.TableModel) (bool, string) {
	b := tm.RequestInCount+tm.ImpressionInCount+tm.ImpressionOutCount+tm.DeliverCount+tm.DeliverBid+tm.Profit > 0
	if b {
		// (supplier,source,time_id,request_in_count,imp_in_count,imp_out_count,deliver_count,deliver_bid,profit)
		return true, fmt.Sprintf(`("%s","%s",%d,%d,%d,%d,%d,%d,%d)`,
			tm.Supplier,
			tm.Source,
			tm.Time,
			tm.RequestInCount,
			tm.ImpressionInCount,
			tm.ImpressionOutCount,
			tm.DeliverCount,
			(tm.DeliverBid)-(tm.Profit),
			tm.Profit,
		)
	}

	return false, ""
}

const supDemSrcTable = `INSERT INTO sup_dem_src
(supplier,demand,source,time_id,request_out_count,imp_out_count,imp_in_count,win_count,win_bid,deliver_count,deliver_bid,profit) VALUES
%s
ON DUPLICATE KEY UPDATE
 request_out_count=request_out_count+VALUES(request_out_count),
 imp_out_count=imp_out_count+VALUES(imp_out_count),
 imp_in_count=imp_in_count+VALUES(imp_in_count),
 win_count=win_count+VALUES(win_count),
 win_bid=win_bid+VALUES(win_bid),
 deliver_count=show_count+VALUES(deliver_count),
 deliver_bid=deliver_bid+VALUES(deliver_bid),
 profit=profit+VALUES(profit)
`
const supSrcTable = `INSERT INTO sup_src
(supplier,source,time_id,request_in_count,imp_in_count,imp_out_count,deliver_count,deliver_bid,profit) VALUES
%s
ON DUPLICATE KEY UPDATE
 request_in_count=request_in_count+VALUES(request_in_count),
 imp_in_count=imp_in_count+VALUES(imp_in_count),
 imp_out_count=imp_out_count+VALUES(imp_out_count),
 deliver_count=deliver_count+VALUES(deliver_count),
 deliver_bid=deliver_bid+VALUES(deliver_bid),
 profit=profit+VALUES(profit)
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
