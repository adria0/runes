package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/adriamb/runes/server/config"
	"github.com/adriamb/runes/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// C is the package config
var C config.Config

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "runes",
	Short: "A personal markdown pad",
	Long:  "A personal markdown pad",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var runCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long:  "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		json, _ := json.MarshalIndent(C, "", "  ")
		fmt.Println("Efective configuration: " + string(json))
		startServer(C)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the repository",
	Long:  "Creates the repository",
	Run: func(cmd *cobra.Command, args []string) {
		store := store.NewStore(C.DataDir)
		err := store.Entry.Create()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Repository initialized")
	},
}

// ExecuteCmd adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteCmd() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)

	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.runes.yaml)")
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(runCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".runes") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.SetEnvPrefix("runes")   // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.Unmarshal(&C); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Configuration file ~/.runes.yaml not found, using default settings.")
		C.WebServer.Bind = ":8086"
		C.Auth.Type = config.AuthNone
	}

	if C.DataDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		C.DataDir = usr.HomeDir + "/.runes"
	}

	if C.TmpDir == "" {
		C.TmpDir = "/tmp/runes/tmp"
	}

	if C.CacheDir == "" {
		C.CacheDir = "/tmp/runes/cache"
	}

}
