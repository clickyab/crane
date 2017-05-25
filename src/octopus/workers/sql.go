package workers

import (
	"bytes"
	"fmt"
	"octopus/models"
	"services/broker"
	"text/template"
)

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
	Acknowledger  *broker.Delivery
}

const sqlTemplate = `{{define "sql1""}}
INSERT INTO sup_dem_src
(supplier,demand,source,time_id,request,impression,win,show_time,imp_bid,show_bid, win_bid) VALUES
{{rows .SupDemSrc "sup_dem_src"}}
ON DUPLICATE KEY UPDATE
 request=request+VALUES(request), impression=impression+VALUES(impression),win=win+VALUES(win), show_time=show_time+VALUES(show_time),imp_bid=imp_bid+VALUES(imp_bid),
 show_bid=show_bid+VALUES(show_bid),win_bid=win_bid+VALUES(win_bid);
{{end}}
{{define "sql2""}}
 INSERT INTO sup_src
(supplier,source,time_id,request,impression,show_time,imp_bid,show_bid) VALUES
{{rows .SupSrc "sup_src"}}
ON DUPLICATE KEY UPDATE
 request=request+VALUES(request), impression=impression+VALUES(impression), show_time=show_time+VALUES(show_time), imp_bid=imp_bid+VALUES(imp_bid), show_bid=show_bid+VALUES(show_bid);
 {{end}}`

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
			c++
		}
	}
	if tableName == "sup_dem_src" {
		c := 1
		for _, r := range m {
			end := ""
			if len(m) != c {
				end = ","
			}
			end += "\n"
			res += fmt.Sprintf("(%s,%s,%d,%d,%d,%d,%d,%d)%s",
				r.Supplier, r.Source, r.Time, r.Request, r.Impression, r.Show, r.ImpressionBid, r.ShowBid, end)
			c++
		}
	}
	return res
}})

func flush(supDemSrc map[string]*tableModel, supSrc map[string]*tableModel) error {
	buf1 := bytes.Buffer{}
	queryTemplate.Lookup("sql1").Execute(&buf1, struct {
		SupDemSrc map[string]*tableModel
	}{
		supDemSrc,
	})
	q1 := buf1.String()
	buf2 := bytes.Buffer{}
	queryTemplate.Lookup("sql2").Execute(&buf1, struct {
		SupDemSrc map[string]*tableModel
	}{
		supDemSrc,
	})
	q2 := buf2.String()
	return models.NewManager().UpdateConsume(q1, q2)
}

func init() {
	queryTemplate.Parse(sqlTemplate)
}
