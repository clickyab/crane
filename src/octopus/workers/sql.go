package workers

import (
	"fmt"
	"text/template"
)

const sql = `INSERT INTO ads  (name, age) VALUES ("ali", 3""),("reza", 4 "")  ON DUPLICATE KEY UPDATE age=age+VALUES(age)`

type tableModel struct {
	Supplier      string // all
	Source        string // all
	Demand        string // all except impression
	Time          int    // all
	Request       int    // impression
	Impression    int    // impression slot
	Win           int    // demand
	Show          int    // show
	ImpressionBid int64  // winner
	ShowBid       int64  // show
	WinnerBid     int64  // Winner
	Acknowledger  *Acknowledger
}

const sqlTemplate = `{{define sql}}
START TRANSACTION;
INSERT INTO sup_dem_src
(supplier,demand,source,time,request,impression,win,show,imp_bid,show_bid, win_bid) VALUES
{{rows .SupDemSrc "sup_dem_src"}}
ON DUPLICATE KEY UPDATE SET
 request=request+VALUES(request), impression=impression+VALUES(impression),win=win+VALUES(win), show=show+VALUES(show),imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid),win_bid=win_bid+VALUES(win_bid);

 INSERT INTO sup_src
(supplier,source,time,request,impression,show,imp_bid,show_bid) VALUES
{{rows .SupSrc "sup_src"}}
ON DUPLICATE KEY UPDATE SET
 request=request+VALUES(request), impression=impression+VALUES(impression), show=show+VALUES(show), imp_bid=imp_bid+VALUES(imp_bid), show_bid=show_bid+VALUES(show_bid);
 `

var queryTemplate = template.New("sql").Funcs(template.FuncMap{"rows": func(m map[string]*tableModel, tableName string) string {
	res := ""
	c := 1
	// IF EDITING THIS BE CAREFUL AND DOUBLE CHECK WITH TEMPLATE TO BE SURE BOTH ARE IN THE SAME ORDER
	if tableName == "sup_src" {
		for _, r := range m {
			end := ""
			if len(m) != c {
				end = ","
			}
			end += "\n"
			res += fmt.Sprintf("(%s,%s,%s,%d,%d,%d,%d,%d,%d,%d,%d)%s",
				r.Supplier, r.Demand, r.Source, r.Time, r.Request, r.Impression, r.Win, r.Show, r.ImpressionBid, r.ShowBid, r.WinnerBid, end)
		}
	}
	if tableName == "sup_dem_src" {
		for _, r := range m {
			end := ""
			if len(m) != c {
				end = ","
			}
			end += "\n"
			res += fmt.Sprintf("(%s,%s,%d,%d,%d,%d,%d,%d)%s",
				r.Supplier, r.Source, r.Time, r.Request, r.Impression, r.Show, r.ImpressionBid, r.ShowBid, end)
		}
	}

	return res
}})

func flush(supDemSrc map[string]*tableModel, supSrc map[string]*tableModel) error {
	// TODO insert into both tables
	return nil
}

func init() {

}
