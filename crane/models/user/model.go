package user

import (
	"time"

	"fmt"

	"clickyab.com/crane/crane/models/publisher"
)

// User user model in database
// @Model {
//		table = users
//		primary = true, id
//		find_by = id
//		transaction = insert
//		list = yes
// }
type User struct {
	ID        int64                  `json:"id" db:"id"`
	Email     string                 `json:"email" db:"email"`
	Domain    string                 `json:"domain" db:"domain"`
	Password  string                 `json:"password" db:"password"`
	Active    publisher.ActiveStatus `json:"active" db:"active"`
	CreatedAt *time.Time             `json:"created_at"  db:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at" db:"updated_at"`
}

func (m *Manager) IsUserActive(ID int64) bool {
	res := User{}
	q := fmt.Sprintf("SELECT * FROM %s WHERE active=?", UserTableFull)
	err := m.GetRDbMap().SelectOne(&res, q)
	if err != nil {
		return false
	}
	return true
}
