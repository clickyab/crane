npm-install:
	cd $(ROOT)/js/banner && $(NPM) install
	cd $(ROOT)/js/vmap && $(NPM) install

build-js: npm-install
	cd $(ROOT)/js/banner && $(NPM) run build
	rm -rf $(ROOT)/dist
	cp -r $(ROOT)/js/banner/dist $(ROOT)/js
	cp $(ROOT)/js/vmap/node_modules/vmap-kit/dist/jwplayer/vmap.js $(ROOT)/js/dist/jwplayer.js
	cp $(ROOT)/js/vmap/node_modules/vmap-kit/dist/videojs/vmap.js $(ROOT)/js/dist/videojs.js

js: build-js go-bindata
	cd $(ROOT)/js/dist && $(BIN)/go-bindata -nometadata -o $(ROOT)/supplier/layer/internal/js/js.gen.go -nomemcopy=true -pkg=js ./
