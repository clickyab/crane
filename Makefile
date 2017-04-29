export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
include $(ROOT)/build/variables.mk
all: $(GB) codegen
	$(BUILD)
include $(ROOT)/build/common.mk
include $(ROOT)/build/gb.mk
include $(ROOT)/build/linter.mk
include $(ROOT)/build/bindata.mk
include $(ROOT)/build/migration.mk
include $(ROOT)/build/codegen.mk
include $(ROOT)/build/services.mk
include $(ROOT)/build/cleanup.mk
include $(ROOT)/build/test.mk
