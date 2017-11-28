package builder

import (
	"math/rand"
	"net"

	"clickyab.com/gad/models"
	"clickyab.com/crane/crane/entity"
)

// ShowOptionSetter is the function to handle setting
type ShowOptionSetter func(*context) (*context, error)


// common is the common type data
type common struct {
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

// App is the common application data
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

	Slots []*Slot
}

// Data is calculated or fetched data
type Data struct {
	Website      *models.Website
	App          *models.App
	PhoneData    *models.PhoneData
	CellLocation *models.CellLocation
}

// context is the app context
type context struct {
	common common
	app    App
	rtb    RTB
	data   Data

	showT int
}

func (c *context) IP() net.IP {
	return c.common.IP
}

func (c *context) OS() entity.OS {
	return c.common.OS
}

func (c *context) ClientID() string {
	panic("implement me")
}

func (c *context) Protocol() string {
	panic("implement me")
}

func (c *context) UserAgent() string {
	panic("implement me")
}

func (c *context) Location() entity.Location {
	panic("implement me")
}

func (c *context) Attributes() map[string]string {
	panic("implement me")
}

func (c *context) TrackID() string {
	panic("implement me")
}

func (c *context) Publisher() entity.Publisher {
	panic("implement me")
}

func (c *context) Slots() []entity.Slot {
	panic("implement me")
}

func (c *context) Category() []entity.Category {
	panic("implement me")
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

// GetCommon return th common part of data
func (c *context) GetCommon() *common {
	return &c.common
}

// GetCommon return th app part of data
func (c *context) GetApp() *App {
	return &c.app
}

// GetCommon return th rtb part of data
func (c *context) GetRTB() *RTB {
	return &c.rtb
}

// GetCommon return th data part of data
func (c *context) GetData() *Data {
	return &c.data
}

// ShowT is a hack to handle a simple redirection (for Clickyab owners need)
func (c *context) ShowT() bool {
	if c.showT == 0 {
		c.showT = 2
		if c.common.Mobile && c.common.ProvinceID > 0 && c.common.Alexa && rand.Intn(chanceShowT.Int()) == 1 {
			c.showT = 3
		}
	}

	return c.showT == 3
}

// GetPublisher return the current publisher object
func (c *context) GetPublisher() entity.Publisher {
	if c.data.App == nil && c.data.Website == nil {
		return invalidPub{}
	}

	if c.data.App == nil {
		return c.data.Website
	}
	return c.data.App
}

// NewContext return a context based on its setters
func NewContext(opt ...ShowOptionSetter) (*context, error) {
	res := &context{}
	var err error
	for i := range opt {
		if res, err = opt[i](res); err != nil {
			return nil, err
		}
	}

	return res, err
}
