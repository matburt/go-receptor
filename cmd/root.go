package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dataDir string
	debug   bool
	peer    []string
	rootCmd = &cobra.Command{
		Use:   "receptor",
		Short: "A mesh worker networking system",
	}
)

// Execute main
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Turn on verbose output")
	rootCmd.PersistentFlags().StringVar(&dataDir,
		"data_dir",
		"/var/lib/receptor",
		"Path to the directory where Receptor stores its database and metadata")
	// TODO: IPSlice or similar?
	rootCmd.PersistentFlags().StringSliceVarP(&peer, "peer", "p", nil, "Peers to connect to directly")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("peer", rootCmd.PersistentFlags().Lookup("peer"))
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".receptor")
		viper.AddConfigPath("/etc/receptor/")
		viper.AddConfigPath("$HOME/.receptor")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Config file not found, skipping")
			} else {
				panic(fmt.Errorf("Fatal error config file: %s", err))
			}
		}
	}
	viper.SetDefault("Debug", false)
	viper.AutomaticEnv()
}
