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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SamWolfs/rofi-github/github"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repositoriesCmd)
}

var repositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "List all user-owned and user-configured repositories.",
	Long: `List all user-owned and user-configured repositories.
By default only user-owned repositories will be returned;
to add repositories owned by other users, manually add them to your configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			var repositories []github.Repository
			mapstructure.Decode(metadata.Get("repositories"), &repositories)
			rows := make([]string, len(repositories))
			for i, repo := range repositories {
				// Provide the url as "hidden" info, available from the $ROFI_INFO environment variable
				rows[i] = formatResource(repo)
			}
			fmt.Println("\x00markup-rows\x1ftrue")
			fmt.Print(strings.Join(rows, "\n"))
		} else {
			response := os.Getenv("ROFI_INFO")
			if response == "" {
				log.Fatal("Choose a repository.")
				return
			}
			var repository github.Repository
			json.Unmarshal([]byte(response), &repository)
			repository.View()
		}
	},
}
