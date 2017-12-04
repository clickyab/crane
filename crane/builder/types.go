package builder

import (
	"math/rand"
	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models"
)

// ShowOptionSetter is the function to handle setting
type ShowOptionSetter func(*Context) (*Context, error)

// Common is the Common type data
type Common struct {
	Type           string
	IP             net.IP
	ISPID          int64
	ProvinceID     int64
	Country        string
	City           string
	Province       string
	Isp            string
	UserAgent      string
	Browser        string
	OS             entity.OS
	Platform       string
	PlatformID     int64
	BrowserVersion string
	Method         string
	Referrer       string
	Mobile         bool
	Host           string
	Scheme         string
	MegaImp        string
	CopID          string
	TID            string
	Parent         string
	Alexa          bool
	NoTiny         bool
}

// App is the Common application data
type App struct {
	//other App stuff
	GoogleID      string
	AndroidID     string
	AndroidDevice string
	SDKVersion    int64
}

// RTB is the real time biding data
type RTB struct {
	MinCPC           int64
	MinBidPercentage float64
	FloorDIV         int64
	MultiVideo       bool
	UnderFloor       bool
	EventPage        string
	Async            bool // Default is sync, each request must return the data
	NoCap            bool // Do not use capping system. default is false

	Slots []entity.Seat
}

// Data is calculated or fetched data
type Data struct {
	Website      *models.Website
	App          *models.App
	PhoneData    *models.PhoneData
	CellLocation *models.CellLocation
}

type invalidPub struct {
}

func (invalidPub) GetID() int64 {
	return 0
}

func (invalidPub) GetName() string {
	return ""
}

func (invalidPub) FloorCPM() int64 {
	return 0
}

func (invalidPub) GetActive() bool {
	return false
}

func (invalidPub) GetType() string {
	return ""
}

// ShowT is a hack to handle a simple redirection (for Clickyab owners need)
func (c *Context) ShowT() bool {
	if c.showT == 0 {
		c.showT = 2
		if c.common.Mobile && c.common.ProvinceID > 0 && c.common.Alexa && rand.Intn(chanceShowT.Int()) == 1 {
			c.showT = 3
		}
	}

	return c.showT == 3
}

// GetPublisher return the current publisher object
func (c *Context) GetPublisher() entity.Publisher {
	if c.data.App == nil && c.data.Website == nil {
		return invalidPub{}
	}

	if c.data.App == nil {
		return c.data.Website
	}
	return c.data.App
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
