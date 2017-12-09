package jwt

import (
	"context"
	"crypto/rsa"

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
	privatePem = config.RegisterString("services.store.jwt.private_256", ``, "")
	publicPem  = config.RegisterString("services.store.jwt.public_256", ``, "")
)

func (j *initJwt) Initialize(ctx context.Context) {
	var err error
	private, err = crypto.ParseRSAPrivateKeyFromPEM([]byte(privatePem.String()))

	assert.Nil(err)

	public, err = crypto.ParseRSAPublicKeyFromPEM([]byte(publicPem.String()))
	assert.Nil(err)
}

func init() {
	initializer.Register(&initJwt{}, 0)
}
