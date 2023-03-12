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
package github

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/cli/go-gh"
)

type WorkflowResponse struct {
	TotalCount int32 `json:"total_count"`
	Workflows  []Workflow
}

type Workflow struct {
	Id         int    `mapstructure:"id"`
	Name       string `mapstructure:"name"`
	Repository string `mapstructure:"repository"`
}

func (w Workflow) Format(color string) string {
	info, _ := json.Marshal(w)
	return fmt.Sprintf("%s <span color=\"%s\" size=\"10pt\" style=\"italic\">(%s)</span>\x00info\x1f%s", w.Name, color, w.Repository, string(info[:]))
}

func (w Workflow) View() {
	_, stdErr, err := gh.Exec("workflow", "view", "-R", w.Repository, "--web", strconv.Itoa(w.Id))
	if err != nil {
		log.Fatal(string(stdErr.Bytes()[:]))
	}
}

func GetUserWorkflows() []Workflow {
	return GetWorkflows(GetUserRepositories())
}

func GetWorkflows(repositories []Repository) []Workflow {
	var workflows []Workflow
	for _, repo := range repositories {
		workflows = append(workflows, getWorkflowsForRepository(repo)...)
	}
	return workflows
}

func getWorkflowsForRepository(repository Repository) []Workflow {
	query := fmt.Sprintf("/repos/%s/%s/actions/workflows", repository.Owner.Login, repository.Name)
	stdOut, _, err := gh.Exec("api", query)
	if err != nil {
		panic(err)
	}

	var workflowResponse WorkflowResponse
	if err := json.Unmarshal(stdOut.Bytes(), &workflowResponse); err != nil {
		panic(err)
	}

	workflows := workflowResponse.Workflows[:0]
	for _, workflow := range workflowResponse.Workflows {
		workflow.Repository = fmt.Sprintf("%s/%s", repository.Owner.Login, repository.Name)
		workflows = append(workflows, workflow)
	}

	return workflows
}
