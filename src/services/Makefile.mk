export SERVICES_ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export GO?=$(shell which go)
export UPDATE?=
export GOPATH?=$(shell mktemp -d)
export SERVICES_LINTER=$(GOPATH)/bin/gometalinter.v1
export SERVICES_LINTERCMD?=$(SERVICES_LINTER) -e ".*.gen.go" -e "upstream.go" -e "tmp" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=misspell

$(SERVICES_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(SERVICES_ROOT)/tmp/ip2l
	wget -c http://www.clickyab.com/downloads/IP-COUNTRY-REGION-CITY.BIN -O $(SERVICES_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN

service_bindata:
	GOPATH=$(SERVICES_ROOT)/tmp $(GO) get -v $(UPDATE) github.com/jteeuwen/go-bindata/go-bindata

$(SERVICES_ROOT)/ip2location/data.gen.go: $(SERVICES_ROOT)/tmp/ip2l/IP-COUNTRY-REGION-CITY.BIN service_bindata
	cd $(SERVICES_ROOT)/tmp/ip2l && $(SERVICES_ROOT)/tmp/bin/go-bindata -nomemcopy -o $(SERVICES_ROOT)/ip2location/data.gen.go -pkg ip2location .

services_dummy_ip2l:
	@[ -f $(SERVICES_ROOT)/ip2location/data.gen.go ] || make -f $(SERVICES_ROOT)/Makefile.mk $(SERVICES_ROOT)/ip2location/data.gen.go

services_ip2l: services_dummy_ip2l

services_convey:
	$(GO) get -v github.com/smartystreets/goconvey/...

services_linter:
	$(GO) get -v gopkg.in/alecthomas/gometalinter.v1
	$(SERVICES_LINTER) --install

services_test: services_ip2l services_convey services_linter
	rm -rf $(GOPATH)/src/services
	mkdir -p $(GOPATH)/src/services && cp -r $(SERVICES_ROOT)/* $(GOPATH)/src/services/
	$(SERVICES_LINTERCMD) $(GOPATH)/src/services/...
	$(GO) get -v github.com/smartystreets/goconvey/...
	cd $(GOPATH)/src/services && $(GO) get -v ./...
	cd $(GOPATH)/src/services && $(GO) test -v ./...
	rm -rf $(GOPATH)/src/services

