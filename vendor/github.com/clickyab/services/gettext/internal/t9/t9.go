package t9

import "fmt"

// Base is the Translation base
type Base struct {
	Text   string        `json:"text"`
	Params []interface{} `json:"params"`
}

// GetText return base text
func (t9 Base) GetText() string {
	return t9.Text
}

// GetParams return all params
func (t9 Base) GetParams() []interface{} {
	return t9.Params
}

func (t9 Base) String() string {
	return fmt.Sprintf(t9.Text, t9.Params...)
}
