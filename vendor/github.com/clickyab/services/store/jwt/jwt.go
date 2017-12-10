package jwt

import (
	"fmt"
	"time"

	"errors"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/store"
)

// NewJWT implement bearer interface. private and public are rsa 256 pem key
func NewJWT() store.Bearer {
	return &enc{}
}

type enc struct {
}

const (
	exp    = "exp"
	format = "060102150405"
)

func (c *enc) Encode(m map[string]string, d time.Duration) string {
	cl := jws.Claims{}

	for k, v := range m {
		cl.Set(k, v)
	}
	cl.Set(exp, time.Now().Add(d).Format(format))
	cl.SetIssuedAt(time.Now())
	t := jws.NewJWT(cl, crypto.SigningMethodRS256)
	res, err := t.Serialize(private)
	assert.Nil(err)
	return string(res)
}

func (c *enc) Decode(b []byte, ks ...string) (bool, map[string]string, error) {
	j, err := jws.ParseJWT(b)
	if err != nil {
		return false, nil, err
	}
	if err := j.Validate(public, crypto.SigningMethodRS256); err != nil {
		return false, nil, err
	}
	if !j.Claims().Has(exp) {
		return false, nil, errors.New("expire time is not set")
	}
	tm := j.Claims().Get(exp).(string)
	tx, err := time.Parse(format, tm)
	if err != nil {
		return false, nil, errors.New("expire time format is not valid")
	}
	isExp := false
	if tx.Before(time.Now()) {
		isExp = true
	}
	res := make(map[string]string)
	if len(ks) != 0 {
		for _, v := range ks {
			if !j.Claims().Has(v) {
				return isExp, nil, fmt.Errorf("key %s not found", v)
			}
			res[v] = j.Claims().Get(v).(string)
		}
	} else {
		for k := range j.Claims() {
			if k != exp && k != "iat" {
				res[k] = fmt.Sprint(j.Claims().Get(k))
			}
		}
	}
	return isExp, res, nil
}
