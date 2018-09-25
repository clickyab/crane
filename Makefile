export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export APPNAME=crane
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export NODE=$(shell which nodejs)
export NPM=$(shell which npm)
export GIT=$(shell which git)
export BIN=$(ROOT)/bin
export GOPATH=$(abspath $(ROOT)/../../..)
export GOBIN?=$(BIN)
export DIFF=$(shell which diff)
export WATCH?=hello
export ORIGIN_GIT_DIR?=$(ROOT)/.git
export LONG_HASH?=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORT_HASH?=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMIT_DATE?=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMP_DATE=$(shell date +%Y%m%d)
export COMMIT_COUNT?=$(shell git rev-list HEAD --count| cat)
export BUILD_DATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X version.hash=$(LONG_HASH) -X version.short=$(SHORT_HASH) -X version.date=$(COMMIT_DATE) -X version.count=$(COMMIT_COUNT) -X version.build=$(BUILD_DATE)"
export LDARG=-ldflags $(FLAGS)
export BUILD=cd $(ROOT) && $(GO) install -v $(LDARG)
export DBPASS?=$(DEFAULT_PASS)
export DB_USER?=root
export DB_NAME?=clickyab
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp
export UGLIFYJS=$(ROOT)/node_modules/.bin/uglifyjs
export CRN_SERVICES_MYSQL_WDSN=root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true
export CRN_SERVICES_MYSQL_RDSN?=dev:cH3M7Z7I4sY8QP&ll130U&73&6KS@tcp(db-2.clickyab.ae:3306)/clickyab?charset=utf8&parseTime=true
export CRN_CRANE_SUPPLIER_STREAM_ADDRESS=127.0.0.1:9801

all: codegen
	$(BUILD) ./...

run-demand: all ip2location
	$(ROOT)/bin/demand

debug-demand: debuger
	cd $(ROOT)/bin && $(BIN)/dlv --listen=:2345 --headless=true --api-version=2 debug clickyab.com/crane/commands/demand

run-supplier: all ip2location
	$(ROOT)/bin/supplier

debug-supplier: debuger
	$(BIN)/dlv --listen=:2345 --headless=true --api-version=2 exec $(BIN)/supplier

run-imp: all ip2location
	$(ROOT)/bin/impression-worker

debug-imp: ip2location debuger
	$(BIN)/dlv --listen=:2345 --headless=true --api-version=2 exec $(BIN)/impression-worker

run-click: all ip2location
	$(ROOT)/bin/click-worker

debug-click: ip2location debuger
	$(BIN)/dlv --listen=:2345 --headless=true --api-version=2 exec $(BIN)/click-worker

include $(ROOT)/scripts/*.mk
