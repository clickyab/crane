convey: $(GB)
	$(BUILD) github.com/smartystreets/goconvey

mockgen: $(GB)
	$(BUILD) github.com/golang/mock/mockgen

mockentity: $(LINTER) mockgen
	mkdir -p $(ROOT)/src/octopus/exchange/mock_exchange
	$(BIN)/mockgen -destination=$(ROOT)/src/octopus/exchange/mock_exchange/mock_exchange.gen.go octopus/exchange Impression,Demand,Advertise,Publisher,Location,Slot,Supplier


.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity

test-gui: mockentity codegen convey
	cd $(ROOT)/src && goconvey -host=0.0.0.0

test: mockentity codegen convey
	cd $(ROOT) && $(GB) test -v