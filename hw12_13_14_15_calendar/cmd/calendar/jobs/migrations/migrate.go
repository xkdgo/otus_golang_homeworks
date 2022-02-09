package migrations

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"

	// need import to implement migrations at this block.
	_ "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/migrations"
)

const serviceName = "calendar"

var (
	FromEnv bool
	CfgFile string
)

func Migrate(cmd *cobra.Command, args []string) error {
	fmt.Println("running migrate ...")
	var err error
	var conf config.CalendarConfig
	fmt.Println("fromenv = ", FromEnv)
	if FromEnv {
		conf, err = config.NewCalendarConfigFromEnv(serviceName)
	} else {
		conf, err = config.NewCalendarConfigFromFile(CfgFile, serviceName)
	}
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("need one more argument <up|down|status>")
	}

	db, err := goose.OpenDBWithDriver("pgx", conf.Storage.DSN)
	if err != nil {
		return errors.Wrapf(err, "goose: failed to open DB: %v", conf.Storage.DSN)
	}

	err = db.Ping()
	if err != nil {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			err = db.Ping()
			if err == nil {
				break
			}
		}
		if err != nil {
			return errors.Wrapf(err, "goose: DB ping failed: %v", conf.Storage.DSN)
		}
	}
	_, err = db.Exec("SELECT 1")
	if err != nil {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			_, err = db.Exec("SELECT 1")
			if err == nil {
				break
			}
		}
		if err != nil {
			return errors.Wrapf(err, "goose: DB select test failed: %v", conf.Storage.DSN)
		}
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	migrateDir := os.Getenv("MIGRATION_DIR")
	if migrateDir == "" {
		migrateDir, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(migrateDir)
	if err = goose.Run(args[0], db, migrateDir); err != nil {
		return errors.Wrapf(err, ":goose migration error with arg %s", args[0])
	}
	return nil
}
