tools-codegen:
	$(BUILD) clickyab.com/exchange/commands/codegen

octopus-user: tools-codegen
	$(BIN)/codegen -p clickyab.com/exchange/octopus/console/user/aaa
	$(BIN)/codegen -p clickyab.com/exchange/octopus/console/user/routes


codegen: ip2location migration octopus-user
