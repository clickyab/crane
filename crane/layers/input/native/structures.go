package native

import (
	"github.com/Sirupsen/logrus"

	"net"

	"context"
	"errors"
	"fmt"
	"strconv"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"clickyab.com/crane/crane/models/publisher"
	"clickyab.com/crane/crane/models/user"
	"github.com/clickyab/services/random"
)

type nativeImp struct {
}

// Transform Transform request to impression
func (n *nativeImp) Transform(ctx context.Context, r entity.Request) (entity.Impression, error) {
	//fetch website by domain and (clickyab) supplier
	pubManager := publisher.NewPublisherManager()
	if r.Attributes()["d"] == "" || r.Attributes()["supplier"] == "" {
		return nil, errors.New("domain or supplier empty")
	}
	pub, err := pubManager.FindPublisherByPlatformNameSup(r.Attributes()["d"], publisher.WebPlatform, r.Attributes()["supplier"])
	if err != nil {
		return nil, errors.New("publisher with the specified domain not found")
	}
	userManager := user.NewUserManager()
	ok := userManager.IsUserActive(pub.UserID)
	if !ok {
		return nil, errors.New("user not found or inactive")
	}
	res := &imp{}
	under := pub.FUnderFloor == publisher.ActiveStatusTrue
	res.fPub = &local.Publisher{
		FName:         pub.FName,
		FSupplier:     pub.FSupplier,
		FFloorCPM:     pub.FFloorCPM,
		FSoftFloorCPM: pub.FSoftFloorCPM,
		FUnderFloor:   &under,
	}
	res.fTrackID = <-random.ID

	intCount, err := strconv.Atoi(r.Attributes()["count"])
	if err != nil {
		return nil, errors.New("count not valid")
	}
	sResult := make([]entity.Slot, intCount)
	attr := make(map[string]interface{})
	for i := 1; i <= intCount; i++ {
		//fill slot
		sResult = append(sResult, local.ExtractSlot(pub.FSupplier, pub.FName, publisher.NativePlatform, 0, 0, fmt.Sprintf("%d", i), attr))
	}
	res.fSlots = sResult
	// TODO implement later
	res.fCategories = make([]entity.Category, 0)
	return res, nil

}

type imp struct {
	fAttr       map[string]string
	fTrackID    string
	fClientID   string
	fIP         net.IP
	fUA         string
	fPub        *local.Publisher
	fLocation   entity.Location
	fOS         entity.OS
	fSlots      []entity.Slot
	fCategories []entity.Category
	fprotocol   string

	nDum   []entity.Slot
	latlon entity.LatLon
}

func (i *imp) Attributes() map[string]string {
	return i.fAttr
}

// JustForLint TODO :// remove it afterwards
func JustForLint(i imp) {
	var _ = nativeImp{}
	if false {
		b := i.latlon
		logrus.Debug(b)
		i.extractData()
	}
	return
}
