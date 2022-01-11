package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/migrations"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate cmd is used for database migration",
		Long:  `migrate cmd is used for database migration: migrate <up | down>`,
		RunE:  Migrate,
	}
)

func init() {
	migrateCmd.Flags().StringVarP(&CfgFile, "config", "c", "/etc/calendar/config.toml", "Path to configuration file")
	viper.BindPFlag("config", migrateCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(migrateCmd)
}

func Migrate(cmd *cobra.Command, args []string) error {
	fmt.Println("running migrate ...")
	config, err := NewConfig(CfgFile)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("need one more argument <up|down>")
	}

	db, err := sql.Open("pgx", config.Storage.DSN)
	if err != nil {
		return err
	}

	if err = goose.Run(args[0], db, "./"); err != nil {
		return err
	}
	return nil
}
