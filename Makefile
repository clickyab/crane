export APPNAME=gad
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export NODE=$(shell which nodejs)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GB=$(BIN)/gb
export LINTER=$(BIN)/gometalinter.v1
export GOPATH=$(ROOT):$(ROOT)/vendor
export DIFF=$(shell which diff)
export WATCH?=hello
export LONGHASH=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMPDATE=$(shell date +%Y%m%d)
export COMMITCOUNT=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X version.hash=$(LONGHASH) -X version.short=$(SHORTHASH) -X version.date=$(COMMITDATE) -X version.count=$(COMMITCOUNT) -X version.build=$(BUILDDATE)"
export LDARG=-ldflags $(FLAGS)
export BUILD=$(BIN)/gb build $(LDARG)
export DBPASS?=$(DEFAULT_PASS)
export DB_USER?=root
export DB_NAME?=clickyab
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp
export LINTERCMD=$(LINTER) -e ".*.gen.go" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell
export UGLIFYJS=$(ROOT)/node_modules/.bin/uglifyjs

.PHONY: all gb clean

all: $(GB)
	$(BUILD)

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

gb:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get -v github.com/constabulary/gb/...

metalinter:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO)  get -v gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(LINTER) --install

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin

$(GB):
	@[ -f $(BIN)/gb ] || make gb

$(LINTER):
	@[ -f $(LINTER) ] || make metalinter

fetch: $(GB)
	PATH=$(PATH):$(BIN) $(GB) vendor fetch $(REPO)

crane: $(GB)
	$(BUILD) cmd/crane

run-crane: crane
	sudo setcap cap_net_bind_service=+ep $(BIN)/crane
	$(BIN)/crane

mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql | true
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql | true
	echo 'FLUSH PRIVILEGES;' | mysql mysql | true
	echo 'DROP DATABASE IF EXISTS $(DB_NAME); CREATE DATABASE $(DB_NAME);' | mysql -u $(DB_USER) -p$(DBPASS)
	#mysql -u $(DB_USER) -p$(DBPASS) -c $(DB_NAME) <$(ROOT)/db/structure.sql

rabbitmq-setup: needroot
	[ "1" -eq "$(shell rabbitmq-plugins enable rabbitmq_management | grep 'Plugin configuration unchanged' | wc -l)" ] || (rabbitmqctl stop_app && rabbitmqctl start_app)
	rabbitmqctl add_user $(RUSER) $(RPASS) || rabbitmqctl change_password $(RUSER) $(RPASS)
	rabbitmqctl set_user_tags $(RUSER) administrator
	rabbitmqctl set_permissions -p / $(RUSER) ".*" ".*" ".*"
	wget -O /usr/bin/rabbitmqadmin http://127.0.0.1:15672/cli/rabbitmqadmin
	chmod a+x /usr/bin/rabbitmqadmin
	rabbitmqadmin declare queue name=dlx-queue
	rabbitmqadmin declare exchange name=dlx-exchange type=topic
	rabbitmqctl set_policy DLX ".*" '{"dead-letter-exchange":"dlx-exchange"}' --apply-to queues
	rabbitmqadmin declare binding source="dlx-exchange" destination_type="queue" destination="dlx-queue" routing_key="#"

SUBDIRS := $(wildcard $(ROOT)/src/*)

lint: $(LINTER) $(SUBDIRS)

$(SUBDIRS):
	$(LINTERCMD) $@/...

uglifyjs:
	npm install uglifyjs

$(UGLIFYJS):
	@[ -f $(UGLIFYJS) ] || make uglifyjs

uglify: $(UGLIFYJS)
	rm -rf $(ROOT)/tmp/embed
	mkdir -p $(ROOT)/tmp/embed
	cp $(ROOT)/statics/show.js $(ROOT)/tmp/embed/show.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/show.js -o $(ROOT)/tmp/embed/show-min.js
	cp $(ROOT)/statics/conversion/clickyab-tracking.js $(ROOT)/tmp/embed/clickyab-tracking.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/conversion/clickyab-tracking.js -o $(ROOT)/tmp/embed/clickyab-tracking-min.js
	cp $(ROOT)/statics/vastAD.js $(ROOT)/tmp/embed/vastAD.js
	$(NODE) $(UGLIFYJS) $(ROOT)/statics/vastAD.js -o $(ROOT)/tmp/embed/vastAD-min.js

go-bindata: $(GB)
	$(BUILD) github.com/jteeuwen/go-bindata/go-bindata

embed: go-bindata uglify
	cd $(ROOT)/tmp/embed/ && $(BIN)/go-bindata -o $(ROOT)/src/statics/static-no-lint.go -nomemcopy=true -pkg=statics ./...

create-imp-table :
	echo 'CREATE TABLE impressions$(IMPDATE)  LIKE impressions20161108; ' | mysql -u $(DB_USER) -p$(DBPASS) -c $(DB_NAME)

restore: $(GB)
	PATH=$(PATH):$(BIN) $(GB) vendor restore
	cp $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done

conditional-restore:
	$(DIFF) $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done || make restore

docker-build: conditional-restore all

ansible:
	ansible-playbook -vvvv -i $(ROOT)/contrib/deploy/hosts.ini $(ROOT)/contrib/deploy/staging.yaml

convey: $(GB)
	$(BUILD) github.com/smartystreets/goconvey

mockgen: $(GB)
	$(BUILD) github.com/golang/mock/mockgen
	mkdir -p $(ROOT)/src/entity/mock_entity
	mkdir -p $(ROOT)/src/clickyab/entity/mock_entity

mockentity: $(LINTER) mockgen
	$(BIN)/mockgen -destination=$(ROOT)/src/entity/mock_entity/mock_entity.gen.go entity Impression,Demand,Advertise,Publisher,Location,Slot,Supplier

test-gui: mockentity convey
	cd $(ROOT)/src && goconvey -host=0.0.0.0

test: mockentity convey
	cd $(ROOT) && $(GB) test -v

.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity