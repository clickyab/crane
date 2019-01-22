package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/clickyab/services/kv"

	"github.com/clickyab/services/random"

	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

// key of context value
type key int

const (
	uidKey     = "a8f5f167f44f4964e6c998dee827110c"
	lsKey      = "1ab482db173cf15b7fdbe7c38feb6748"
	passphrase = "7fc0cf07fc5f0f57dba7fb538c3b7532"
	userPrefix = "USR"
	// KEY of user in context
	KEY = key(1000)
)

type middleware struct {
}

func (middleware) PreRoute() bool {
	return true
}

func extractList(c string, u *openrtb.User) {
	rls, err := decrypt([]byte(c), passphrase)
	if err != nil {
		return
	}
	res := &openrtb.UserData{
		Name:    "list",
		Id:      "1",
		Segment: []*openrtb.UserData_Segment{},
	}

	for _, v := range strings.Split(string(rls), ",") {
		k := kv.NewEavStore(fmt.Sprintf("%s_%s_%s", userPrefix, u.GetId(), v))
		if len(k.AllKeys()) == 0 {
			continue
		}
		xl := make([]string, len(k.AllKeys()))
		for lk := range k.AllKeys() {
			xl = append(xl, lk)
		}

		res.Segment = append(res.Segment, &openrtb.UserData_Segment{
			Id:    v,
			Value: strings.Join(xl, ","),
		})
	}
}

func (middleware) Handler(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// domain := r.URL.Host
		// fmt.Println("DOMAIN:", domain)
		// if parts := strings.Split(r.URL.Hostname(), "."); len(parts) > 2 {
		// 	domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
		// }
		user := &openrtb.User{
			Data: make([]*openrtb.UserData, 0),
		}
		if uc, err := r.Cookie(uidKey); err != nil {
			user.Id = <-random.ID
			http.SetCookie(w,
				&http.Cookie{
					// Domain:  "." + domain,
					Expires: time.Now().AddDate(2, 0, 0),
					Value:   user.Id,
					Name:    uidKey,
					Path:    "/",
				})
		} else {
			user.Id = uc.Value
		}

		if ls, err := r.Cookie(lsKey); err != nil {
			ec, err := encrypt([]byte(""), passphrase)
			if err == nil {
				http.SetCookie(w,
					&http.Cookie{
						// Domain:  "." + domain,
						Expires: time.Now().AddDate(2, 0, 0),
						Value:   string(ec),
						Name:    lsKey,
						Path:    "/",
					})
			}
		} else {
			extractList(ls.Value, user)

		}
		next(context.WithValue(ctx, KEY, user), w, r)
	}
}

func init() {
	router.RegisterGlobalMiddleware(&middleware{})
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err

	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err

	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err

	}
	return plaintext, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
