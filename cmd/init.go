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

	"github.com/SamWolfs/rofi-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

// TODO: Improve logging/code
var initCmd = &cobra.Command{
	Use: "init",
	Short: "Initialize the configuration files.",
	Long: `Initialize the configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if metadata.IsSet("repositories") {
			fmt.Println("Repositories already defined.")
		} else {
			repos := github.GetUserRepositories()
			repositories := ReadRepositories(repos)
			metadata.Set("repositories", repositories)
			fmt.Println("Set default repositories to user-owned.")
		}
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
