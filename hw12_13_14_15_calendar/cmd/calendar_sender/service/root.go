package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"
)

const (
	serviceName         = "sender"
	serviceNameExchange = "scheduler"
)

// rootCmd represents the base command when called without any subcommands.
var (
	CfgFile   string
	FromEnv   bool
	Release   string
	BuildDate string
	GitHash   string
	rootCmd   = &cobra.Command{
		Use:   "calendar_sender",
		Short: "This is the calendar_sender application",
		Long: `This is the calendar_sender application
it can use several storages
all configurations you can find in config file example`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var conf config.SenderConfig
			if FromEnv {
				conf, err = config.NewSenderConfigFromEnv(serviceName)
			} else {
				conf, err = config.NewSenderConfigFromFile(CfgFile, serviceName)
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
	// will be global for your application.
	rootCmd.Flags().BoolVarP(&FromEnv, "fromenv", "e", false, "Configuration parameters from env and some defaults")
	viper.BindPFlag("fromenv", rootCmd.Flags().Lookup("fromenv"))
	rootCmd.Flags().StringVarP(&CfgFile, "config", "c", "/etc/calendar/config_sender.toml", "Path to configuration file")
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))

	rootCmd.AddCommand(cmdVersion)
}
