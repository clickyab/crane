npm-install:
	cd $(ROOT)/js/banner && $(NPM) install

build-js: npm-install
	cd $(ROOT)/js/banner && $(NPM) run build

show-js: build-js go-bindata
	cd $(ROOT)/js/banner/dist && $(BIN)/go-bindata -nometadata -o $(ROOT)/supplier/layer/web/show-js.gen.go -nomemcopy=true -pkg=web ./
