/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/dalesearle/asciitable"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile         string
	includeAll      bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "rabbitcli",
		Short: "CLI tool for managing taxhawk rabbit clusters",
		Long: `
rabbitcli is a proprietary tool used to manage connections and queues for Taxhawk
RabbitMQ clusters. The primary use is to correct multiple logins and repair
connections that are missing rpc queues.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rabbitcli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rabbitcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rabbitcli")
	}

	viper.AutomaticEnv() // read in set variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil && !initialize{
		fmt.Fprintln(os.Stderr, "Config file needs to be initialized, run rabbitcli config -i")
		os.Exit(1)
	}
}

func newAsciiTable() *asciitable.Table {
	table := asciitable.New()
	table.SetHeaderJustification(asciitable.JustifyCenter)
	table.SetDataJustification(asciitable.JustifyLeft)
	table.SetCellPadding(1, 1)
	return table
}
