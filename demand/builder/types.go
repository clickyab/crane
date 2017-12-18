package builder

// ShowOptionSetter is the function to handle setting
type ShowOptionSetter func(*Context) (*Context, error)

type user string

func (u user) ID() string {
	return string(u)
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
