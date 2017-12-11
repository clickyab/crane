go-bindata:
	GOBIN=$(BIN) $(GO) get -v github.com/jteeuwen/go-bindata/go-bindata
	GOBIN=$(BIN) $(GO) install -v github.com/jteeuwen/go-bindata/go-bindata

$(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN:
	mkdir -p $(ROOT)/contrib
	cd $(ROOT)/contrib && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY-ISP.BIN.gz
	cd $(ROOT)/contrib && gunzip IP-COUNTRY-REGION-CITY-ISP.BIN.gz
	cd $(ROOT)/contrib && rm -f IP-COUNTRY-REGION-CITY-ISP.BIN.md5 && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY-ISP.BIN.md5
	cd $(ROOT)/contrib && md5sum -c IP-COUNTRY-REGION-CITY-ISP.BIN.md5

$(BIN)/IP-COUNTRY-REGION-CITY-ISP.BIN: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN
	cp $(ROOT)/contrib/IP-COUNTRY-REGION-CITY-ISP.BIN $(BIN)

ip2location: $(BIN)/IP-COUNTRY-REGION-CITY-ISP.BIN

prepare: ip2location

.PHONY: go-bindata