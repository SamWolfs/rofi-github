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

	// "github.com/SamWolfs/rofi-github/github"
	"github.com/cli/go-gh"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repositoriesCmd)
}

var repositoriesCmd = &cobra.Command{
	Use: "repositories",
	Short: "List all user-owned and user-configured repositories.",
	Long: `List all user-owned and user-configured repositories.
By default only user-owned repositories will be returned;
to add repositories owned by other users, add them to your configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var repositories []Repository
		mapstructure.Decode(metadata.Get("repositories"), &repositories)
		if len(args) == 0 {
			names := make([]string, len(repositories))
			for i, repo := range repositories {
				names[i] = repo.Name
			}
			fmt.Print(strings.Join(names, "\n"))
		} else {
			choice := args[0]
			url := ""
			for _, repo := range repositories {
				if repo.Name == choice {
					url = repo.Url
					break;
				}
			}
			if url == "" {
				log.Fatal("Choose a repository.")
				return
			}
			_, stdErr, err := gh.Exec("repo", "view", "--web", url)
			if err != nil {
				log.Fatal(stdErr.Bytes()[:])
			}
		}
	},
}
