npm-install:
	cd $(ROOT)/js && $(NPM) install

build-js: npm-install
	cd $(ROOT)/js && $(NPM) run build

show-js: build-js go-bindata
	cd $(ROOT)/js/dist && $(BIN)/go-bindata -nometadata -o $(ROOT)/supplier/layer/web/show-js.gen.go -nomemcopy=true -pkg=web ./
