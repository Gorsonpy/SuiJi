include ./common.mk

SERVICE_NAME = runFzu

.PHONY: build
build:
	sh build.sh

.PHONY: new
new:
	hz new \
	-module $(MODULE) \
	hz update -idl ./idl/api.thrift

.PHONY: gen
gen:
	hz update -idl ./idl/api.thrift
	hz update -idl ./idl/goal.thrift
	hz update -idl ./idl/multiLedger.thrift
	hz update -idl ./idl/ledger.thrift

.PHONY: server
server:
	make build
	cd output && sh bootstrap.sh