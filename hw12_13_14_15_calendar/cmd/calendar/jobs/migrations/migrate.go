package migrations

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"
)

var CfgFile string

func Migrate(cmd *cobra.Command, args []string) error {
	fmt.Println("running migrate ...")
	conf, err := config.NewConfig(CfgFile)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("need one more argument <up|down|status>")
	}

	db, err := sql.Open("pgx", conf.Storage.DSN)
	if err != nil {
		return errors.Wrapf(err, ":migration load driver error")
	}

	if err = goose.Run(args[0], db, "./"); err != nil {
		return errors.Wrapf(err, ":goose migration error with arg %s", args[0])
	}
	return nil
}
