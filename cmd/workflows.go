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

	"github.com/cli/go-gh"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workflowsCommand)
}

var workflowsCommand = &cobra.Command{
	Use: "workflows",
	// TODO: descriptions
	Short: "List all user-owned and user-configured workflows.",
	Long: `List all user-owned and user-configured workflows.
By default only user-owned workflows will be returned;
to add workflows owned by other users, add them to your configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			var workflows []Workflow
			mapstructure.Decode(metadata.Get("workflows"), &workflows)
			rows := make([]string, len(workflows))
			for i, workflow := range workflows {
				info, _ := json.Marshal(workflow)
				rows[i] = workflow.Name + "\x00info\x1f" + string(info[:])
			}
			fmt.Print(strings.Join(rows, "\n"))
		} else {
			response := os.Getenv("ROFI_INFO")
			if response == "" {
				log.Fatal("Choose a workflow.")
				return
			}
			var workflow Workflow
			json.Unmarshal([]byte(response), &workflow)
			_, stdErr, err := gh.Exec("workflow", "view", "-R", workflow.Repository, "--web", workflow.Id)
			if err != nil {
				log.Fatal(string(stdErr.Bytes()[:]))
			}
		}
	},
}
