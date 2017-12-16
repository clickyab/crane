export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export APPNAME=crane
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export NODE=$(shell which nodejs)
export GIT=$(shell which git)
export BIN=$(ROOT)/bin
export GOPATH=$(abspath $(ROOT)/../../..)
export GOBIN?=$(BIN)
export DIFF=$(shell which diff)
export WATCH?=hello
export ORIGIN_GIT_DIR?=$(ROOT)/.git
export LONGHASH?=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH?=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE?=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMPDATE=$(shell date +%Y%m%d)
export COMMITCOUNT?=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X version.hash=$(LONGHASH) -X version.short=$(SHORTHASH) -X version.date=$(COMMITDATE) -X version.count=$(COMMITCOUNT) -X version.build=$(BUILDDATE)"
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
export CRN_SERVICES_MYSQL_RDSN?=root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true

all: codegen
	$(BUILD) ./...

run-webserver: all ip2location
	$(ROOT)/bin/crane

run-imp: all ip2location
	$(ROOT)/bin/crane-impression-worker

include $(ROOT)/scripts/*.mk
