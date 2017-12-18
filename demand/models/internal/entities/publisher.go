package entities

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clickyab.com/crane/demand/entity"
)

// FindWebsiteByPublicID return publisher id for public id
func FindWebsiteByPublicID(pid int64) (entity.Publisher, error) {
	ws := make([]Website, 0)
	_, err := NewManager().GetRDbMap().Select(&ws, `SELECT * from websites where w_pub_id=?`, pid)
	if err != nil {
		return nil, err
	}
	if len(ws) != 1 {
		panic(fmt.Sprintf("there is more then one record with public id %d", pid))
	}
	return &ws[0], nil
}

// ErrorNotAllowCreate rise when supplier doesn't allow to add new website or app
var ErrorNotAllowCreate = errors.New("insert not allowed")

// FindOrAddWebsite return publisher id for given supplier,domain
func FindOrAddWebsite(sup entity.Supplier, domain string, pid int64) (int64, error) {

	if pid == 0 {
		pid = PublicIDGen(sup.Name(), domain)
	}

	w, err := FindWebsiteByPublicID(pid)
	if err == nil {
		if w.Supplier().Name() != sup.Name() {
			return 0, fmt.Errorf("mismatch supplier for domain %s with public id %d. suppliers %s and %s",
				domain, pid, sup.Name(), w.Supplier().Name())
		}
		return w.ID(), nil
	}

	ws := make([]Website, 0)
	_, err = NewManager().GetRDbMap().Select(&ws, `SELECT * from websites where w_supplier=? and w_domain=?`,
		sup.Name(), domain)

	if err == nil && len(ws) != 0 {
		var tw = ws[0]

		for i := range ws {
			if tw.totalImp() < ws[i].totalImp() {
				tw = ws[i]
			}
		}
		return tw.ID(), nil
	}

	if err != nil && !sup.AllowCreate() {
		return 0, ErrorNotAllowCreate
	}

	q := `INSERT INTO websites (u_id, w_domain,w_supplier,w_status,created_at,updated_at,w_date, w_pub_id)
VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE
  u_id=VALUES(u_id),w_domain=VALUES(w_domain),w_supplier=VALUES(w_supplier),w_status=VALUES(w_status),
  created_at=VALUES(created_at),updated_at=VALUES(updated_at),w_date=VALUES(w_date), w_id=LAST_INSERT_ID(w_id)`

	t := time.Now()
	res, err := NewManager().GetWDbMap().Exec(q, sup.UserID(),
		sql.NullString{Valid: true, String: domain},
		sup.Name(), 1, sql.NullString{Valid: true, String: t.String()}, sql.NullString{Valid: true, String: t.String()},
		int(t.Unix()), PublicIDGen(sup.Name(), domain),
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()

}
