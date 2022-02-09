package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/cmd/calendar/jobs/migrations"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"
)

const serviceName = "calendar"

// rootCmd represents the base command when called without any subcommands.
var (
	CfgFile   string
	FromEnv   bool
	Release   string
	BuildDate string
	GitHash   string
	rootCmd   = &cobra.Command{
		Use:   "calendar",
		Short: "This is the calendar application",
		Long: `This is the calendar application
it can use several storages
all configurations you can find in config file example`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("fromenv = ", FromEnv)
			var err error
			var conf config.CalendarConfig
			if FromEnv {
				conf, err = config.NewCalendarConfigFromEnv(serviceName)
			} else {
				conf, err = config.NewCalendarConfigFromFile(CfgFile, serviceName)
			}
			if err != nil {
				return err
			}
			fmt.Println(CfgFile)
			fmt.Printf("%#v\n", conf)
			RunApp(conf)
			return nil
		},
	}

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print application version",
		Long:  `Print application version"`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := json.NewEncoder(os.Stdout).Encode(struct {
				Release   string
				BuildDate string
				GitHash   string
			}{
				Release:   Release,
				BuildDate: BuildDate,
				GitHash:   GitHash,
			}); err != nil {
				fmt.Printf("error while decode version info: %v\n", err)
			}
		},
	}

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate cmd is used for database migration",
		Long:  `migrate cmd is used for database migration: migrate <up | down>`,
		RunE:  migrations.Migrate,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application
	rootCmd.Flags().BoolVarP(&FromEnv, "fromenv", "e", false, "Configuration parameters from env and some defaults")
	viper.BindPFlag("fromenv", rootCmd.Flags().Lookup("fromenv"))
	rootCmd.Flags().StringVarP(&CfgFile, "config", "c", "/etc/calendar/config.toml", "Path to configuration file")
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))

	rootCmd.AddCommand(cmdVersion)
	migrateCmd.Flags().BoolVarP(&migrations.FromEnv,
		"fromenv", "e", false, "Configuration parameters from env and some defaults")
	viper.BindPFlag("fromenv", migrateCmd.Flags().Lookup("fromenv"))
	migrateCmd.Flags().StringVarP(
		&migrations.CfgFile,
		"config",
		"c",
		"/etc/calendar/config.toml",
		"Path to configuration file")
	viper.BindPFlag("config", migrateCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(migrateCmd)
}
