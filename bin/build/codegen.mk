tools-codegen:
	$(BUILD) clickyab.com/crane/commands/codegen

octopus-user: tools-codegen
	$(BIN)/codegen -p clickyab.com/crane/crane/models/internal/ad
	$(BIN)/codegen -p clickyab.com/crane/crane/models/internal/campaign
	$(BIN)/codegen -p clickyab.com/crane/crane/models/internal/user
	$(BIN)/codegen -p clickyab.com/crane/crane/models/internal/publisher

codegen: ip2location migration octopus-user
