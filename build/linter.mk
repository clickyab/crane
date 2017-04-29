SUBDIRS := $(wildcard $(ROOT)/src/*)

$(SUBDIRS):
	$(LINTERCMD) $@/...
lint: codegen $(LINTER) $(SUBDIRS)

metalinter:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO)  get -v gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make metalinter