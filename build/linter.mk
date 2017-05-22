export LINTER=$(BIN)/gometalinter.v1
# TODO : Ignoring services/codegen is a bad thing. try to get it back to lint
export LINTERCMD=$(LINTER) -e ".*.gen.go" -e "upstream.go" -e ".*_test.go" -e "services/codegen" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell

SUBDIRS := $(wildcard $(ROOT)/src/*)

$(SUBDIRS):
	$(LINTERCMD) $@/...
lint: codegen $(LINTER) $(SUBDIRS)

metalinter:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO)  get -v gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make metalinter