tools-migrate: migration
	$(BUILD) clickyab.com/crane/commands/migration

migup: tools-migrate
	$(BIN)/migration -action=up -app=crane

migdown: tools-migrate
	$(BIN)/migration -action=down -app=crane

migdown-all: tools-migrate
	$(BIN)/migration -action=down-all -app=crane

migredo: tools-migrate
	$(BIN)/migration -action=redo -app=crane

miglist: tools-migrate
	$(BIN)/migration -action=list -app=crane

migcreate:
	@/bin/bash $(BIN)/create_migration.sh

migration: go-bindata
	cd $(ROOT) && $(BIN)/go-bindata -nometadata -o ./commands/migration/migration.gen.go -nomemcopy=true -pkg=main ./db/...
