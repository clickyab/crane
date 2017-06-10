tools-codegen:
	$(BUILD) clickyab.com/exchange/commands/codegen

octopus-user: tools-codegen
	$(BIN)/codegen -p clickyab.com/exchange/octopus/console/internal/aaa


codegen: $(ROOT)/services/ip2location/data.gen.go migration octopus-user