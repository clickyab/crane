package aaa

import (
	"fmt"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
	"golang.org/x/crypto/bcrypt"
)

const (
	// DemandUserType demand type
	DemandUserType UserType = "demand"

	// SupplierUserType supplier type
	SupplierUserType UserType = "supplier"

	// AdminUserType admin type
	AdminUserType UserType = "admin"

	// DefaultOrder order
	DefaultOrder = "ASC"
)

type (
	//UserType type user
	//@Enum{
	//}
	UserType string
)

// User user model in database
// @Model {
//		table = users
//		primary = true, id
//		find_by = id,token,email
//		transaction = insert
//		list = yes
// }
type User struct {
	ID       int64    `json:"id" db:"id"`
	Email    string   `json:"email" db:"email"`
	Password string   `json:"password" db:"password"`
	Token    string   `json:"token" db:"token"`
	UserType UserType `json:"user_type" db:"user_type"`
}

// GetUserByToken returns user by its token
func (m *Manager) GetUserByToken(token string) (*User, error) {
	query := `SELECT * FROM %s WHERE token=?`
	query = fmt.Sprintf(query, UserTableFull)
	holder := &User{}
	err := m.GetRDbMap().SelectOne(holder, query, token)
	if err != nil {
		return nil, err
	}
	return holder, nil
}

// RegisterUser try to register user
func (m *Manager) RegisterUser(email, pass string, userType UserType) (*User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	assert.Nil(err)
	err = m.Begin()
	assert.Nil(err)
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			assert.Nil(m.Commit())
		}
	}()
	u := &User{
		Email:    email,
		Password: string(password),
		UserType: userType,
		Token:    <-random.ID,
	}
	err = m.CreateUser(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
