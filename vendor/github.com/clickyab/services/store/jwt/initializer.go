package jwt

import (
	"context"
	"crypto/rsa"

	"io/ioutil"

	"github.com/SermoDigital/jose/crypto"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
)

type initJwt struct {
}

var (
	private    *rsa.PrivateKey
	public     *rsa.PublicKey
	privatePem = config.RegisterString("services.store.jwt.private_256", ``, "private key for crypt data")
	publicPem  = config.RegisterString("services.store.jwt.public_256", ``, "public file to crypt data")
	filePem    = config.RegisterBoolean("services.store.jwt.file", false, "the private/public pair are file name")
)

func (j *initJwt) Initialize(ctx context.Context) {
	var err error
	var priv, pub []byte
	if filePem.Bool() {
		pub, err = ioutil.ReadFile(publicPem.String())
		assert.Nil(err)
		priv, err = ioutil.ReadFile(privatePem.String())
		assert.Nil(err)
	} else {
		pub = []byte(publicPem.String())
		priv = []byte(privatePem.String())
	}

	private, err = crypto.ParseRSAPrivateKeyFromPEM(priv)
	assert.Nil(err)

	public, err = crypto.ParseRSAPublicKeyFromPEM(pub)
	assert.Nil(err)

}

func init() {
	initializer.Register(&initJwt{}, 0)
}
