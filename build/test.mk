convey: $(GB)
	$(BUILD) github.com/smartystreets/goconvey

mockgen: $(GB)
	$(BUILD) github.com/golang/mock/mockgen

mockentity: $(LINTER) mockgen
	mkdir -p $(ROOT)/src/octopus/exchange/mock_exchange
	mkdir -p $(ROOT)/src/crane/entity/mock_entity
	$(BIN)/mockgen -destination=$(ROOT)/src/octopus/exchange/mock_exchange/mock_exchange.gen.go octopus/exchange Impression,Demand,Advertise,Publisher,Location,Slot,Supplier
	$(BIN)/mockgen -destination=$(ROOT)/src/crane/entity/mock_entity/mock_entity.gen.go crane/entity Impression,Demand,Advertise,Campaign,Publisher,Location,Slot,Supplier



.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity

test-gui: mockentity codegen convey
	cd $(ROOT)/src && goconvey -host=0.0.0.0

test: mockentity codegen convey
	cd $(ROOT) && $(GB) test -v -race