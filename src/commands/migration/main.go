package main

import (
	"config"
	"flag"
	"fmt"
	"os"
	"services/initializer"
	"text/tabwriter"

	"services/mysql"

	"github.com/Sirupsen/logrus"
	"github.com/rubenv/sql-migrate"
)

var (
	action = flag.String("action", "up", "up/down is supported, default is up")
	n      int
)

// doMigration is my try to migrate on demand. but I don't know if there is more than
// one ins is in memory
func doMigration(dir migrate.MigrationDirection, max int) error {
	// OR: Use migrations from bindata:
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "db/migrations",
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
	config.Initialize()

	var err error

	defer initializer.Initialize()()

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
