package main

import (
	"commands"
	"flag"
	"fmt"
	"os"
	"services/config"
	"services/initializer"
	"text/tabwriter"

	"services/mysql"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/rubenv/sql-migrate"
)

var (
	validApps = []string{"octopus"}

	action = flag.String("action", "up", "up/down is supported, default is up")
	app    = flag.String("app", "", "application to handle migrations valids are "+strings.Join(validApps, ","))
	n      int
)

func validApp(in string) bool {
	for i := range validApps {
		if validApps[i] == in {
			return true
		}
	}

	return false
}

// doMigration is my try to migrate on demand. but I don't know if there is more than
// one ins is in memory
func doMigration(dir migrate.MigrationDirection, max int) error {
	// OR: Use migrations from bindata:
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      fmt.Sprintf("db/%s_migrations", *app),
	}

	var err error
	m := mysql.Manager{}
	if max == 0 {
		n, err = migrate.Exec(m.GetWSQLDB(), "mysql", migrations, dir)
	} else {
		n, err = migrate.ExecMax(m.GetWSQLDB(), "mysql", migrations, dir, max)
	}
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())

	var err error

	defer initializer.Initialize()()

	if !validApp(*app) {
		logrus.Fatalf("app is invalid, valids are %s", strings.Join(validApps, ","))
	}

	if *action == "up" {
		err = doMigration(migrate.Up, 0)
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "down" {
		err = doMigration(migrate.Down, 1)
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "down-all" {
		err = doMigration(migrate.Down, 0)
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "redo" {
		err = doMigration(migrate.Down, 1)
		if err == nil {
			err = doMigration(migrate.Up, 1)
		}
		fmt.Printf("\n\n%d migration is applied\n", n)

	} else if *action == "list" {
		var mig []*migrate.MigrationRecord
		m := mysql.Manager{}
		mig, err = migrate.GetMigrationRecords(m.GetRSQLDB(), "mysql")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "|ID\t|Applied at\t|")
		for i := range mig {
			fmt.Fprintf(w, "|%s\t|%s\t|\n", mig[i].Id, mig[i].AppliedAt)
		}
		_ = w.Flush()
	}

	if err != nil {
		logrus.Panic(err)
	}
}
