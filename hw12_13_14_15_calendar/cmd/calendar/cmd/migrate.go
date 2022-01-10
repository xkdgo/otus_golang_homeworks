package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	fmt.Println(CfgFile)
	_ = config
	if len(args) == 0 {
		return fmt.Errorf("need one more argument <up|down>")
	}
	fmt.Println(args[0])
	return nil
}
