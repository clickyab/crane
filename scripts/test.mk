convey:
	$(GO) get -v github.com/smartystreets/goconvey
	$(GO) install -v github.com/smartystreets/goconvey

mockgen:
	$(GO) get -v github.com/golang/mock/mockgen
	$(GO) install -v github.com/golang/mock/mockgen

mockentity: $(LINTER) mockgen
	mkdir -p $(ROOT)/demand/entity/mock_entity
	rm -rf $(ROOT)/demand/entity/mock_entity/*.gen.go
	$(BIN)/mockgen -destination=$(ROOT)/demand/entity/mock_entity/mock_entity.gen.go clickyab.com/crane/demand/entity Context,Creative,Campaign,Publisher,SelectedCreative,Location,Seat,Supplier,Request,User


.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity

test-gui: mockentity codegen convey
	cd $(ROOT) && goconvey -host=0.0.0.0

test: mockentity codegen convey
	cd $(ROOT) && $(GO) test -v -race ./...
