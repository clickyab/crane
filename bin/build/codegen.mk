tools-codegen:
	$(BUILD) clickyab.com/crane/commands/codegen


octopus-user: tools-codegen

codegen: ip2location migration octopus-user