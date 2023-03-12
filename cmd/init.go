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
	"fmt"
	"log"
	"strings"

	"github.com/SamWolfs/rofi-github/github"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var forceInit string

// TODO: Improve logging/code
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration files.",
	Long:  `Initialize the configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		repositories := initRepositories()
		initWorkflows(repositories)
		if viper.IsSet("updateEvery") {
			fmt.Println("Update delay already defined.")
		} else {
			viper.Set("updateEvery", 86400)
			fmt.Println("Set default update delay to 'daily'.")
		}
		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
		if err := metadata.WriteConfig(); err != nil {
			log.Fatal(err)
		}
	},
}

func initRepositories() []github.Repository {
	var repositories []github.Repository
	if strings.Contains(forceInit, "repositories") || metadata.IsSet("repositories") {
		fmt.Println("Repositories already defined. Re-run the command with the `--force,-f` flag to force initialization.")
		mapstructure.Decode(metadata.Get("repositories"), &repositories)
	} else {
		fmt.Println("Fetching user repositories...")
		repositories = github.GetUserRepositories()
		metadata.Set("repositories", repositories)
		fmt.Println("Set default repositories to user-owned.")
	}
	return repositories
}

func initWorkflows(repositories []github.Repository) []github.Workflow {
	var workflows []github.Workflow
	if strings.Contains(forceInit, "workflows") || metadata.IsSet("workflows") {
		fmt.Println("Workflows already defined. Re-run the command with the `--force,-f` flag to force initialization.")
		mapstructure.Decode(metadata.Get("workflows"), &workflows)
	} else {
		fmt.Println("Fetching user workflows...")
		workflows := github.GetWorkflows(repositories)
		metadata.Set("workflows", workflows)
		fmt.Println("Set default workflows to user-owned.")
	}
	return workflows
}
