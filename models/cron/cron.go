package cron

import (
	"fmt"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/slack"
	"github.com/sirupsen/logrus"
)

type impRes struct {
	Imps      int64  `json:"imps" db:"imps"`
	CPM       int64  `json:"cpm" db:"cpm"`
	Supplier  string `json:"supplier" db:"supplier"`
	Publisher string `json:"publisher" db:"publisher"`
}

type clickRes struct {
	Clicks    int64  `json:"clicks" db:"clicks"`
	CPC       int64  `json:"cpc" db:"cpc"`
	Supplier  string `json:"supplier" db:"supplier"`
	Publisher string `json:"publisher" db:"publisher"`
}

// WebImp select from impression type web
func WebImp(date int) error {
	logrus.Debug("start web impression worker")
	defer func() {
		logrus.Debug("web impression worker done")
	}()
	var webImps []impRes
	q := fmt.Sprintf(`SELECT
		COUNT(imp_id) AS imps,
		COALESCE(SUM(imp_cpm),0) AS cpm,
		s_name AS supplier,
		websites.w_domain AS publisher
		FROM impressions%d AS i
		INNER JOIN websites ON websites.w_id=i.w_id
		WHERE i.w_id IS NOT NULL
		AND i.w_id!=0
		AND i.imp_status=0
		GROUP BY i.w_id`, date)
	_, err := entities.NewManager().GetRDbMap().Select(&webImps, q)
	if err != nil {
		return err
	}
	logrus.Debugf("total web impression: %d", len(webImps))
	wr := entities.NewManager().GetWDbMap()
	for i := range webImps {
		q := "INSERT INTO daily_report (supplier,type,publisher,imps,cpm,date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE imps=VALUES(imps),cpm=VALUES(cpm)"
		_, err := wr.Exec(q, webImps[i].Supplier, "web", webImps[i].Publisher, webImps[i].Imps, webImps[i].CPM, date)
		if err != nil {
			// TODO :// just for debugging
			go func() {
				slack.AddCustomSlack(fmt.Errorf("[WTF] insert in daily report failed publisher %s", webImps[i].Publisher))
			}()
			return err
		}
	}
	return nil
}

// AppImp select from impression type app
func AppImp(date int) error {
	logrus.Debug("start app impression worker")
	defer func() {
		logrus.Debug("web app worker done")
	}()
	var appImps []impRes
	q := fmt.Sprintf(`SELECT
		COUNT(imp_id) AS imps,
		COALESCE(SUM(imp_cpm),0) AS cpm,
		s_name AS supplier,
		apps.app_package AS publisher
		FROM impressions%d AS i
		INNER JOIN apps ON apps.app_id=i.app_id
		WHERE i.app_id IS NOT NULL
		AND i.app_id!=0
		AND i.imp_status=0
		GROUP BY i.app_id`, date)
	_, err := entities.NewManager().GetRDbMap().Select(&appImps, q)
	if err != nil {
		return err
	}
	logrus.Debugf("total app impression: %d", len(appImps))
	wr := entities.NewManager().GetWDbMap()

	for i := range appImps {
		q := "INSERT INTO daily_report (supplier,type,publisher,imps,cpm,date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE imps=VALUES(imps),cpm=VALUES(cpm)"
		_, err := wr.Exec(q, appImps[i].Supplier, "app", appImps[i].Publisher, appImps[i].Imps, appImps[i].CPM, date)
		if err != nil {
			// TODO :// just for debugging
			go func() {
				slack.AddCustomSlack(fmt.Errorf("[WTF] insert in daily report failed publisher %s", appImps[i].Publisher))
			}()
			return err
		}
	}
	return nil
}

// WebClick select from click type web
func WebClick(date int) error {
	logrus.Debug("start web click worker")
	defer func() {
		logrus.Debug("web click worker done")
	}()
	var webClicks []clickRes
	q := fmt.Sprintf(`SELECT COUNT(c_id) AS clicks,
		COALESCE(SUM(c_winnerbid),0) AS cpc,
		websites.w_domain AS publisher,
		c_supplier AS supplier
		FROM clicks INNER JOIN websites ON websites.w_id=clicks.w_id
		WHERE clicks.w_id IS NOT NULL
		AND clicks.w_id!=0
		AND c_status=0
		AND c_date=?
		GROUP BY clicks.w_id`)
	_, err := entities.NewManager().GetRDbMap().Select(&webClicks, q, date)
	if err != nil {
		return err
	}
	logrus.Debugf("total web clicks: %d", len(webClicks))
	wr := entities.NewManager().GetWDbMap()

	for i := range webClicks {
		q := "INSERT INTO daily_report (supplier,type,publisher,clicks,cpc,date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE clicks=VALUES(clicks),cpc=VALUES(cpc)"
		_, err := wr.Exec(q, webClicks[i].Supplier, "web", webClicks[i].Publisher, webClicks[i].Clicks, webClicks[i].CPC, date)
		if err != nil {
			// TODO :// just for debugging
			go func() {
				slack.AddCustomSlack(fmt.Errorf("[WTF] insert in daily report failed publisher %s", webClicks[i].Publisher))
			}()
			return err
		}
	}

	return nil
}

// AppClick select from click type app
func AppClick(date int) error {
	logrus.Debug("start app click worker")
	defer func() {
		logrus.Debug("app click worker done")
	}()
	var appClicks []clickRes
	q := fmt.Sprintf(`SELECT COUNT(c_id) AS clicks,
		COALESCE(SUM(c_winnerbid),0) AS cpc,
		apps.app_package AS publisher,
		c_supplier AS supplier
		FROM clicks INNER JOIN apps ON apps.app_id=clicks.app_id
		WHERE clicks.app_id IS NOT NULL
		AND clicks.app_id!=0
		AND c_status=0
		AND c_date=?
		GROUP BY clicks.app_id`)
	_, err := entities.NewManager().GetRDbMap().Select(&appClicks, q, date)
	if err != nil {
		return err
	}
	logrus.Debugf("total app click: %d", len(appClicks))
	wr := entities.NewManager().GetWDbMap()
	for i := range appClicks {
		q := "INSERT INTO daily_report (supplier,type,publisher,clicks,cpc,date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE clicks=VALUES(clicks),cpc=VALUES(cpc)"
		_, err := wr.Exec(q, appClicks[i].Supplier, "app", appClicks[i].Publisher, appClicks[i].Clicks, appClicks[i].CPC, date)
		if err != nil {
			// TODO :// just for debugging
			go func() {
				slack.AddCustomSlack(fmt.Errorf("[WTF] insert in daily report failed publisher %s", appClicks[i].Publisher))
			}()
			return err
		}
	}
	return nil
}
