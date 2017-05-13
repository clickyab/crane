go-bindata: $(GB)
	$(BUILD) github.com/jteeuwen/go-bindata/go-bindata

$(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN:
	mkdir -p $(ROOT)/contrib
	cd $(ROOT)/contrib && wget -c http://static.clickyab.com/IP-COUNTRY-REGION-CITY.BIN

$(ROOT)/src/services/ip2location/data.gen.go: $(ROOT)/contrib/IP-COUNTRY-REGION-CITY.BIN go-bindata
	[ -f $(ROOT)/src/services/ip2location/data.gen.go ] || (cd $(ROOT)/contrib/ && $(BIN)/go-bindata -nomemcopy -o $(ROOT)/src/services/ip2location/data.gen.go -pkg ip2location .)

ip2location: $(ROOT)/src/services/ip2location/data.gen.go
