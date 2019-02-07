package builder

import (
	"strings"

	openrtb "clickyab.com/crane/openrtb/v2.5"
)

// ShowOptionSetter is the function to handle setting
type ShowOptionSetter func(*Context) (*Context, error)

type user struct {
	id       string
	userdata []*openrtb.UserData
}

func (u user) List() map[string][]string {
	ls := make(map[string][]string)
	for _, v := range u.userdata {
		if v.Name == "list" {
			for _, w := range v.Segment {
				ls[w.Id] = strings.Split(w.Value, ",")
			}
		}
	}
	return ls
}

func (u user) ID() string {
	return u.id
}

// NewContext return a Context based on its setters
func NewContext(opt ...ShowOptionSetter) (*Context, error) {
	res := &Context{}
	var err error
	for i := range opt {
		if res, err = opt[i](res); err != nil {
			return nil, err
		}
	}

	return res, err
}
