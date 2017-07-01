tools-codegen:
	$(BUILD) clickyab.com/crane/commands/codegen

octopus-user: tools-codegen
	$(BIN)/codegen -p clickyab.com/crane/crane/models/ad
	$(BIN)/codegen -p clickyab.com/crane/crane/models/campaign
	$(BIN)/codegen -p clickyab.com/crane/crane/models/user
	$(BIN)/codegen -p clickyab.com/crane/crane/models/publisher

codegen: ip2location migration octopus-user
