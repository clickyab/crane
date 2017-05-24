go-bindata:
	$(GO) get -v github.com/jteeuwen/go-bindata/go-bindata
	$(GO) install -v github.com/jteeuwen/go-bindata/go-bindata

$(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(ROOT)/contrib
	cd $(ROOT)/contrib && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY.BIN.gz
	cd $(ROOT)/contrib && gunzip IP-COUNTRY-REGION-CITY.BIN.gz
	cd $(ROOT)/contrib && rm -f IP-COUNTRY-REGION-CITY.BIN.md5 && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY.BIN.md5
	cd $(ROOT)/contrib && md5sum -c IP-COUNTRY-REGION-CITY.BIN.md5

$(ROOT)/services/ip2location/data.gen.go: go-bindata $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN
	[ -f $(ROOT)/services/ip2location/data.gen.go ] || (cd $(ROOT)/contrib/ && $(BIN)/go-bindata -nomemcopy -o $(ROOT)/services/ip2location/data.gen.go -pkg ip2location .)

ip2location: $(ROOT)/services/ip2location/data.gen.go

prepare: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN

.PHONY: go-bindata