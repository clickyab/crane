package native

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"errors"
	"strconv"
	"clickyab.com/crane/crane/models/publisher"
	"context"
	"github.com/rs/xmux"
	"clickyab.com/crane/crane/models/user"
	"fmt"
	"github.com/clickyab/services/random"
)

type cornerType string
type posType string

const(
	sharpCornerType cornerType ="sharp"
	roundCornerType cornerType ="round"

	topPosType posType ="top"
	bottomPosType posType ="bottom"
	rightPosType posType ="right"
	leftPosType posType ="left"
)

type nativeImp struct {

}

func (*nativeImp) Transform(ctx context.Context,userTrackID string,r *http.Request) (entity.Impression,error) {
	supplier:=xmux.Param(ctx,"supplier")
	if supplier!=""{
		return nil,errors.New("supplier not found")
	}
	rawCount:=r.URL.Query().Get("count")
	if rawCount==""{
		return nil,errors.New("count not found")
	}
	intCount,err:=strconv.ParseInt(rawCount,10,0)
	if err!=nil|| intCount < 1 || intCount > 12{
		return nil,errors.New("cant parse count into integer")
	}
	rawDomain:=r.URL.Query().Get("d")
	if rawDomain==""{
		return nil,errors.New("domain not found")
	}

	corners:=r.URL.Query().Get("corners")
	if corners=="" || (corners!=string(sharpCornerType) && corners!=string(roundCornerType)) {
		corners=string(sharpCornerType)
	}

	title:=r.URL.Query().Get("title")
	if title=="" {
		title="undefined title"
	}

	more:=r.URL.Query().Get("more")

	position:=r.URL.Query().Get("position")
	if position=="" || (position!=string(topPosType) && position!=string(bottomPosType) && position!=string(leftPosType) && position!=string(rightPosType)) {
		position=string(topPosType)
	}

	intFontSize:=13 //default TODO read from config
	rawFontSize:=r.URL.Query().Get("fontSize")
	if rawFontSize!=""{
		a,err:=strconv.ParseInt(rawFontSize,10,0)
		if err==nil{
			intFontSize=int(a)
		}
	}

	//fetch website by domain and (clickyab) supplier
	pubMananger:=publisher.NewPublisherManager()
	pub,err:=pubMananger.FindPublisherByPlatformNameSup(rawDomain,publisher.WebPlatform,supplier)
	if err!=nil{
		return nil,errors.New("publisher with the specified domain not found")
	}
	userManager:=user.NewUserManager()
	ok:=userManager.IsUserActive(pub.UserID)
	if !ok{
		return nil,errors.New("user not found or inactive")
	}

	res:=imp{}
	// fill attributes
	res.FAttr["count"]=intCount
	res.FAttr["corners"]=corners
	res.FAttr["title"]=title
	res.FAttr["more"]=more
	res.FAttr["position"]=position
	res.FAttr["fontSize"]=intFontSize

	res.FRequest=r
	res.FIP=local.IP(r.RemoteAddr)
	res.FUA=r.UserAgent()
	res.FLocation=local.FLocation(res.FIP)

	res.FPub=&local.Publisher{
		FName:pub.FName,
		FSupplier:pub.FSupplier,
		FFloorCPM:pub.FFloorCPM,
		FSoftFloorCPM:pub.FSoftFloorCPM,
		FUnderFloor:&(pub.FUnderFloor==publisher.ActiveStatusTrue),
	}
	res.FTrackID=<-random.ID

	attr:=make(map[string]interface{})

	sResult:=make([]entity.Slot,intCount)

	for i:=1; i<=int(intCount);  i++{
		//fill slot
		sResult=append(sResult,local.ExtractSlot(pub.FSupplier,pub.FName,publisher.NativePlatform,0,0,fmt.Sprintf("%d",i),attr))
	}

	return res,nil

}

type imp struct {
	FRequest    *http.Request          `json:"request"`
	FTrackID    string                 `json:"track_id"`
	FClientID   string                 `json:"client_id"`
	FIP         net.IP                 `json:"ip"`
	FUA         string                 `json:"user_agent"`
	FPub        *local.Publisher       `json:"pub"`
	FLocation   entity.Location        `json:"location"`
	FOS         entity.OS              `json:"os"`
	FSlots      []*local.Slot          `json:"slots"`
	FCategories []entity.Category      `json:"categories"`
	FAttr       map[string]interface{} `json:"attr"`

	nDum   []entity.Slot
	latlon entity.LatLon
}













// JustForLint TODO :// remove it afterwards
func JustForLint(i imp) {
	if false {
		b := i.latlon
		logrus.Debug(b)
		i.extractData()
	}
	return
}
