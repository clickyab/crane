tools-codegen:
	$(BUILD) clickyab.com/exchange/commands/codegen

octopus-user: tools-codegen
	$(BIN)/codegen -p clickyab.com/exchange/octopus/console/internal/aaa
	$(BIN)/codegen -p clickyab.com/exchange/octopus/console/internal/routes


codegen: ip2location migration octopus-user
