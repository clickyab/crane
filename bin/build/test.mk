convey:
	$(GO) get -v github.com/smartystreets/goconvey
	$(GO) install -v github.com/smartystreets/goconvey

mockgen:
	$(GO) get -v github.com/golang/mock/mockgen
	$(GO) install -v github.com/golang/mock/mockgen

mockentity: $(LINTER) mockgen
	mkdir -p $(ROOT)/octopus/exchange/mock_exchange
	mkdir -p $(ROOT)/crane/entity/mock_entity
	$(BIN)/mockgen -destination=$(ROOT)/octopus/exchange/mock_exchange/mock_exchange.gen.go clickyab.com/exchange/octopus/exchange Impression,Demand,Advertise,Publisher,Location,Slot,Supplier
	$(BIN)/mockgen -destination=$(ROOT)/crane/entity/mock_entity/mock_entity.gen.go clickyab.com/exchange/crane/entity Impression,Demand,Advertise,Campaign,Publisher,Location,Slot,Supplier


.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity

test-gui: mockentity codegen convey
	cd $(ROOT) && goconvey -host=0.0.0.0

test: mockentity codegen convey
	cd $(ROOT) && $(GO) test -v -race ./...
