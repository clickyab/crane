go-bindata:
	GOBIN=$(BIN) $(GO) get -v github.com/jteeuwen/go-bindata/go-bindata
	GOBIN=$(BIN) $(GO) install -v github.com/jteeuwen/go-bindata/go-bindata

$(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(ROOT)/contrib
	cd $(ROOT)/contrib && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY.BIN.gz
	cd $(ROOT)/contrib && gunzip IP-COUNTRY-REGION-CITY.BIN.gz
	cd $(ROOT)/contrib && rm -f IP-COUNTRY-REGION-CITY.BIN.md5 && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY.BIN.md5
	cd $(ROOT)/contrib && md5sum -c IP-COUNTRY-REGION-CITY.BIN.md5

ip2location: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN
	cp $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN $(BIN)

prepare: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN

.PHONY: go-bindata