clean:
	rm -rf $(ROOT)/tmp/*
	cd $(ROOT) && git clean -fX ./bin

.PHONY: all clean