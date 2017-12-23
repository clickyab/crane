define TEST_SERVICES_STORE_JWT_PRIVATE_256
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsKPiNmGjixdjOBICajBVDlBZ+A18ETiGFOPMkEM4fhOlXvlM
0jiVAtm73dy+Z9bii2eDXKWDlWudhmpqaJ307N4J+Q2CtN4Q805UmNpq7g74Iptz
TF5iyi8tYt9nmZmO3uOsSlJpo+9VCU6Z7cKS2drdJTQOdOeJkvDV2F3bwL/zePJO
xhqZgNilaF8e5Thhvb4cS27AAi1buuJSXRmyIC1XdbOB7EhFwcsGxyx15jnXY8yn
AwlAwsZeHNZnCWZCbgO8PlGjQy/EVwSoRltbRKiNMP2ulbu9iB1qdOWYukwrHuiJ
xRBu2tnM91X8a4S676x+VAoA7dbA0PiWtxRlKQIDAQABAoIBAQCMQQRME8jslxxk
GACs2kWfAPP+/o4FinEEQ0BZR9aiXO0Q9TgnL2A6DDKcXjsdbkUhVYa7WHybdwB0
CykEem9QaJlYlH61KCIjXo3TdJI1BdPGftHU0Jj2WvFZsXOsRX5owjQ6KyfQUCeg
JTYZ0EYUDzFK6gOUlYfqEapqi1QCvdT04puN2JiKxqmRC5c99JriN/HT60aFfyIc
0YscMm7KYCBgunNwuvVYFyVwSHU/rilkD4/GnaE0eHWzLrebsuWklDfm3uUnuh3g
J5Hrrzm3NxI1292KL2w+xF5/a/hnlMIvspJR4fRl/zqnUgXjQvKu7monDiPOBJ0j
TBIknbDVAoGBAN9ccX/suSDm+af+aLlWJotO1z0olfdWHrjdnCP+D/OVJr50JFRp
6PyEPf/oQeS8SNSEqcaGUTt7maxa65ZkjQTB+XzTIMV7OLk7P6Rtu4LnLsxFEPpK
8LDOdoR5NYX9LPE1mSfmkJUQ8O122Snuj1Ex2pQP7XMN2lWQdz30FitDAoGBAMpz
sKYSzmS+aualxh/Zx3RVMqmSi5fzIco+7Ky6F4JkcSMiPaobQsbbTfXOvH4azotY
MUb9mFCYo7Z8qY2NuqEu0Vw8tjOUlCWoeOJbfUS3detN1+/QWZom+gcvuEen6vT0
0rXW0tKrmhuA1qXSZxSD7QlL41DMMUdu13s+mmkjAoGAGXudg1/Jm0Re1VjEL+jx
StF64wJfYbGUnbjC2KRiw4KPwgjUIEFZEH4x6KFh6yvED2L3T6wBVoz6clyJJkZg
hFtxEFmyEnXFefSVhTfzvbKMd0el8Thdj35urRx7C8dfukFIMPIwKqGdKyp1CkCq
XzQV4PCUQZ0h/MMXp5zSbqsCgYEAht/OJsXOpMVLGMAGH0AKJhGAgAI2Z5O9xixn
pqcPdHoP+ZUdOc+RjEOcS49gei7HvYOSyLW4HLGkF7Yziy+Jz0oOhoGX7QMmM3Rq
nHrRGM+Uip/ApW7L0uv2lIURIwPWfzz/h89Hgrx6HaqW1cA3li5R42igVzrB9dH9
Uokhe7sCgYEAw4Yyuf8kHCmvp5eVZHnXr1xnFikUIpc7qa1dZYkdKUVj/ZSOt28J
30dzGftjtnfaBJAeBSzirWJvHJ5vW29UtTK2C7/IpioNj5oacjak3sZrIjmRK7zi
9r8vOe2jxnnNAgjsnAWuOAdeNYKtV34pByb42sx/mKrDvwEGEgvnIN8=
-----END RSA PRIVATE KEY-----
endef
define TEST_SERVICES_STORE_JWT_PUBLIC_256
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsKPiNmGjixdjOBICajBV
DlBZ+A18ETiGFOPMkEM4fhOlXvlM0jiVAtm73dy+Z9bii2eDXKWDlWudhmpqaJ30
7N4J+Q2CtN4Q805UmNpq7g74IptzTF5iyi8tYt9nmZmO3uOsSlJpo+9VCU6Z7cKS
2drdJTQOdOeJkvDV2F3bwL/zePJOxhqZgNilaF8e5Thhvb4cS27AAi1buuJSXRmy
IC1XdbOB7EhFwcsGxyx15jnXY8ynAwlAwsZeHNZnCWZCbgO8PlGjQy/EVwSoRltb
RKiNMP2ulbu9iB1qdOWYukwrHuiJxRBu2tnM91X8a4S676x+VAoA7dbA0PiWtxRl
KQIDAQAB
-----END PUBLIC KEY-----
endef

export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO=$(shell which go)
export GIT=$(shell which git)
export BIN=$(ROOT)/bin
export GOPATH=$(abspath $(ROOT)/../../../..)
export GOBIN?=$(BIN)
export DIFF=$(shell which diff)
export LINTER=$(BIN)/gometalinter.v1
export TEST_SERVICES_STORE_JWT_PRIVATE_256
export TEST_SERVICES_STORE_JWT_PUBLIC_256

# TODO : Ignoring services/codegen is a bad thing. try to get it back to lint
export LINTERCMD=$(LINTER) -e ".*.gen.go" -e ".*_test.go" -e "codegen/.*" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell
metalinter:
	$(GO) get -v gopkg.in/alecthomas/gometalinter.v1
	$(GO) install -v gopkg.in/alecthomas/gometalinter.v1
	$(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make -f $(ROOT)/Makefile metalinter

dependencies:
	cd $(ROOT) && $(GO) get -t -v ./...

test: dependencies
	cd $(ROOT) && $(GO) test ./...

all: dependencies
	cd $(ROOT) && $(GO) build ./...

lint: dependencies $(LINTER)
	cd $(ROOT) && $(LINTERCMD) $(ROOT)/...
