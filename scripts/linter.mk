export LINTER=$(BIN)/gometalinter
export LINTERCMD=$(LINTER) -e ".*.gen.go" -e ".*_test.go" -e "clickyab.com/crane/vendor/.*" --cyclo-over=19  --sort=path --disable-all --line-length=120 --deadline=100s --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=megacheck --enable=misspell

lint: codegen $(LINTER)
	$(LINTERCMD) $(ROOT)/commands/...
	$(LINTERCMD) $(ROOT)/crane/...

metalinter:
	$(GO) get -v github.com/alecthomas/gometalinter
	$(GO) install -v github.com/alecthomas/gometalinter
	$(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make -f $(ROOT)/Makefile metalinter