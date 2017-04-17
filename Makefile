export APPNAME=exchange
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export NODE=$(shell which nodejs)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
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
export DB_NAME?=$(APPNAME)
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp
export LINTERCMD=$(LINTER) -e ".*.gen.go" --cyclo-over=19 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell
export UGLIFYJS=$(ROOT)/node_modules/.bin/uglifyjs

all: $(GB) codegen
	$(BUILD)

include $(ROOT)/src/services/Makefile.mk

.PHONY: all gb clean


needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

gb:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get -v github.com/constabulary/gb/...

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin

$(GB):
	@[ -f $(BIN)/gb ] || make gb

mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql | true
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql | true
	echo 'FLUSH PRIVILEGES;' | mysql mysql | true
	echo 'DROP DATABASE IF EXISTS $(DB_NAME); CREATE DATABASE $(DB_NAME);' | mysql -u $(DB_USER) -p$(DBPASS)

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


$(SUBDIRS):
	echo $@ | grep services  || $(LINTERCMD) $@/...

codegen: services_ip2l migration

go-bindata: $(GB)
	$(BUILD) github.com/jteeuwen/go-bindata/go-bindata

restore: $(GB)
	PATH=$(PATH):$(BIN) $(GB) vendor restore
	cp $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done

conditional-restore:
	$(DIFF) $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done || make restore


migup: tools-migrate
	$(BIN)/migration -action=up

migdown: tools-migrate
	$(BIN)/migration -action=down

migdown-all: tools-migrate
	$(BIN)/migration -action=down-all

migredo: tools-migrate
	$(BIN)/migration -action=redo

miglist: tools-migrate
	$(BIN)/migration -action=list

migcreate:
	@/bin/bash $(BIN)/create_migration.sh

migration: go-bindata
	cd $(ROOT) && $(BIN)/go-bindata -o ./src/commands/migration/migration.gen.go -nomemcopy=true -pkg=main ./db/migrations/...

tools-migrate: $(BIN)/gb migration
	$(BUILD) commands/migration


# Tests ans lint
lint: codegen $(LINTER) $(SUBDIRS)

metalinter:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO)  get -v gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(LINTER) --install

$(LINTER):
	@[ -f $(LINTER) ] || make metalinter

convey: $(GB)
	$(BUILD) github.com/smartystreets/goconvey

mockgen: $(GB)
	$(BUILD) github.com/golang/mock/mockgen
	mkdir -p $(ROOT)/src/entity/mock_entity

mockentity: $(LINTER) mockgen
	$(BIN)/mockgen -destination=$(ROOT)/src/entity/mock_entity/mock_entity.gen.go entity Impression,Demand,Advertise,Publisher,Location,Slot,Supplier

test-gui: mockentity codegen convey
	cd $(ROOT)/src && goconvey -host=0.0.0.0

test: mockentity codegen convey
	cd $(ROOT) && $(GB) test -v

# END Tests ans lint

.PHONY: lint $(SUBDIRS) $(ENTITIES) mockentity