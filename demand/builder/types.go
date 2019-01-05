package builder

// ShowOptionSetter is the function to handle setting
type ShowOptionSetter func(*Context) (*Context, error)

type user struct {
	id    string
	cyuId string
}

func (u user) CyuId() string {
	return u.cyuId
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
