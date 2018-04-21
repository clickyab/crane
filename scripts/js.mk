npm-install:
	cd $(ROOT)/js/banner && $(NPM) install
	cd $(ROOT)/js/tracker && $(NPM) install
	rm -r $(ROOT)/js/tracker/node_modules/@types/lodash
	rm -rf $(ROOT)/js/vmap/node_modules
	rm -rf $(ROOT)/js/vmap/package-lock.json
	cd $(ROOT)/js/vmap && $(NPM) install
	cd $(ROOT)/js/native && $(NPM) install

build-js: npm-install
	cd $(ROOT)/js/banner && $(NPM) run build
	cd $(ROOT)/js/tracker && $(NPM) run build
	rm -rf $(ROOT)/js/supplier && mkdir $(ROOT)/js/supplier
	rm -rf $(ROOT)/js/demand && mkdir $(ROOT)/js/demand
	cp -r $(ROOT)/js/banner/dist/* $(ROOT)/js/supplier
	cp $(ROOT)/js/tracker/dist/index.js $(ROOT)/js/demand/tracker.js
	cp $(ROOT)/js/vmap/node_modules/vmap-kit/dist/jwplayer/vast.js $(ROOT)/js/supplier/jwplayer.js
	sed -i 's/registerPlugin("vast/registerPlugin("{{.PLUG}}/' $(ROOT)/js/supplier/jwplayer.js
	cp $(ROOT)/js/vmap/node_modules/vmap-kit/dist/videojs/vast.js $(ROOT)/js/supplier/videojs.js
	sed -i 's/videojs.plugin("vast/videojs.plugin("clickyabAdScheduler/' $(ROOT)/js/supplier/videojs.js
	cd $(ROOT)/js/native && $(NPM) run build
	cp $(ROOT)/js/native/build/browser/index.cjs.js $(ROOT)/js/supplier/native.js

js: build-js go-bindata
	cd $(ROOT)/js/supplier && $(BIN)/go-bindata -nometadata -o $(ROOT)/supplier/layers/internal/js/js.gen.go -nomemcopy=true -pkg=js ./
	cd $(ROOT)/js/demand && $(BIN)/go-bindata -nometadata -o $(ROOT)/demand/layers/internal/js/js.gen.go -nomemcopy=true -pkg=js ./
