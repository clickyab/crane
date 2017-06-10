package aaa

import "fmt"

const usersTable string = "users"

// User user model in database
// @Model {
//		table = users
//		primary = true, id
//		find_by = id,token,email
//		transaction = insert
//		list = yes
// }
type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Token    string `json:"token" db:"token"`
}

// GetUserByToken returns user by its token
func (m *Manager) GetUserByToken(token string) (*User, error) {
	query := `SELECT * FROM %s WHERE token=?`
	query = fmt.Sprintf(query, usersTable)
	holder := &User{}
	err := m.GetRDbMap().SelectOne(holder, query, token)
	if err != nil {
		return nil, err
	}
	return holder, nil
}
