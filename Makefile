export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
include $(ROOT)/bin/build/variables.mk
all: codegen
	$(BUILD) ./...
include $(ROOT)/bin/build/common.mk
include $(ROOT)/bin/build/gb.mk
include $(ROOT)/bin/build/linter.mk
include $(ROOT)/bin/build/bindata.mk
include $(ROOT)/bin/build/migration.mk
include $(ROOT)/bin/build/codegen.mk
include $(ROOT)/bin/build/services.mk
include $(ROOT)/bin/build/cleanup.mk
include $(ROOT)/bin/build/test.mk
