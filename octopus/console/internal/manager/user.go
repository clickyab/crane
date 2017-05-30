package manager

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

const usersTable string = "users"

// User user model in database
type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Token    string `json:"token" db:"token"`
	Email    string `json:"email" db:"email"`
}

// GetUserByToken returns user by its token
func (m *Manager) GetUserByToken(token string) (*User, error) {
	query := `SELECT * FROM %s WHERE token=?`
	query = fmt.Sprintf(query, usersTable)
	holder := &User{}
	err := NewManager().GetRDbMap().SelectOne(holder, query, token)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}
	return holder, nil
}
