gb:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get -v github.com/constabulary/gb/...


$(GB):
	@[ -f $(BIN)/gb ] || make gb


restore: $(GB)
	PATH=$(PATH):$(BIN) $(GB) vendor restore
	cp $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done

conditional-restore:
	$(DIFF) $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done || make restore