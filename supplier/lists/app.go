package lists

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/clickyab/services/kv"

	openrtb "clickyab.com/crane/openrtb/v2.5"
)

const prefix = "LST"

func genKey(s string) kv.Kiwi {
	return kv.NewEavStore(fmt.Sprintf("%s_%s", prefix, s))

}

func getData(s string) (kv.Kiwi, *openrtb.UserData) {

	res := &openrtb.UserData{
		Id:      "1",
		Name:    "list",
		Segment: []*openrtb.UserData_Segment{},
	}
	kiwi := genKey(s)

	if len(kiwi.AllKeys()) == 0 {
		return kiwi, res
	}

	seq := make(map[string][]string)

	for k, v := range kiwi.AllKeys() {
		if _, ok := seq[k]; ok {
			seq[k] = append(seq[k], v)
		} else {
			seq[k] = []string{v}
		}
	}

	for k, v := range seq {
		res.Segment = append(res.Segment, &openrtb.UserData_Segment{
			Id:    k,
			Name:  k,
			Value: strings.Join(v, ","),
		})
	}

	return kiwi, res
}

// GetLists return user list
func GetLists(ctx context.Context, uid string) (*openrtb.UserData, error) {
	k, u := getData(uid)
	_ = k.Save(time.Hour * 24 * 30)
	return u, nil

}

// SetLists will update the user list
func SetLists(ctx context.Context, uid, url string, lid ...string) error {
	k, _ := getData(uid)
	for i := range lid {
		ut := append(strings.Split(k.SubKey(lid[i]), ","), url)
		if len(ut) > 50 {
			ut = ut[len(ut)-50:]
		}
		k.SetSubKey(lid[i], strings.Join(ut, ","))
	}
	return k.Save(time.Hour * 24 * 30)
}

// SetListsClean will update the user list
func SetListsClean(ctx context.Context, uid, url string, cls []string, lid ...string) error {
	k, _ := getData(uid)
	for i := range lid {
		k.SetSubKey(lid[i], url)
	}
	_ = k.Drop()
	nk := genKey(uid)

	for i := range cls {
		if ts := k.SubKey(cls[i]); ts != "" {
			nk.SetSubKey(cls[i], ts)
		}
	}

	for i := range lid {
		ut := append(strings.Split(k.SubKey(lid[i]), ","), url)
		if len(ut) > 50 {
			ut = ut[len(ut)-50:]
		}
		nk.SetSubKey(lid[i], strings.Join(ut, ","))
	}

	return k.Save(time.Hour * 24 * 30)
}
