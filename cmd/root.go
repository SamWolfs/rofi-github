/*
Copyright © 2023 Sam Wolfs be.samwolfs@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var metadata = viper.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rofi-github",
	Short: "Interact with your GitHub resources through Rofi.",
	Long: `Rofi GitHub is a collection of CLI utilities to enable developers
to interact with their GitHub resources (repositories, workflows, ...),
using Rofi as the interface.

Provide a resource name as the first argument to get started.
e.g. rofi-github workflows`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is rofi-github/config in the UserConfigDir as defined in: https://pkg.go.dev/os#UserConfigDir)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	initCmd.Flags().StringVarP(&forceInit, "force", "f", "", "Force resource to be re-initialized")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory.
		userConfigDir, err := os.UserConfigDir()
		configDir := userConfigDir + "/rofi-github"

		createConfigIfNotExists(configDir, "config")
		createConfigIfNotExists(configDir, "metadata")
		cobra.CheckErr(err)

		// Search config in config directory with name "config" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// Search metadata in config directory with name "metadata" (without extension).
		metadata.AddConfigPath(configDir)
		metadata.SetConfigType("yaml")
		metadata.SetConfigName("metadata")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	if err := metadata.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using metadata file:", metadata.ConfigFileUsed())
	}
}

func createConfigIfNotExists(configDir string, fileName string) {
	fullPath := configDir + "/" + fileName
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(configDir, os.ModePerm)
		file, err := os.Create(fullPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}
