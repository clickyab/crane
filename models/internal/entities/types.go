package entities

import (
	"database/sql"
	"database/sql/driver"
	"strings"
)

type size struct {
	Width,
	Height int32
}

var (
	sizes = map[int32]*size{
		1:  {Width: 120, Height: 600},
		2:  {Width: 160, Height: 600},
		3:  {Width: 300, Height: 250},
		4:  {Width: 336, Height: 280},
		5:  {Width: 468, Height: 60},
		6:  {Width: 728, Height: 90},
		7:  {Width: 120, Height: 240},
		8:  {Width: 320, Height: 50},
		9:  {Width: 800, Height: 440},
		11: {Width: 300, Height: 600},
		12: {Width: 970, Height: 90},
		13: {Width: 970, Height: 250},
		14: {Width: 250, Height: 250},
		15: {Width: 300, Height: 1050},
		16: {Width: 320, Height: 480},
		17: {Width: 48, Height: 320},
		18: {Width: 128, Height: 128},
	}
)

// SharpArray type is the hack to handle # splited text in our database
type SharpArray []string

// Scan convert the json array ino string slice
func (pa *SharpArray) Scan(src interface{}) error {
	s := &sql.NullString{}
	err := s.Scan(src)
	if err != nil {
		return err
	}

	if s.Valid {
		var res []string
		for _, v := range strings.Split(s.String, "#") {
			if strings.Trim(v, "\n\t ") != "" {
				res = append(res, v)
			}
		}
		*pa = SharpArray(res)
	} else {
		*pa = []string{}
	}
	return nil

}

// Value try to get the string slice representation in database
func (pa SharpArray) Value() (driver.Value, error) {
	s := sql.NullString{}
	s.Valid = len(pa) != 0
	s.String = "#" + strings.Join(pa, "#") + "#"

	return s.Value()
}

// Array is the function to get array of string of this
func (pa SharpArray) Array() []string {
	return []string(pa)
}
