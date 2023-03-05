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
	"github.com/cli/go-gh"
)

type WorkflowResponse struct {
	TotalCount int32 `json:"total_count"`
	Workflows  []Workflow
}

type Workflow struct {
	Id         int
	Name       string
	Repository string
}

func GetUserWorkflows() []Workflow {
	var workflows []Workflow
	for _, repo := range GetUserRepositories() {
		workflows = append(workflows, getWorkflowsForRepository(repo)...)
	}
	return workflows
}

func getWorkflowsForRepository(repository Repository) []Workflow {
	query := "/repos/" + repository.Owner.Login + "/" + repository.Name + "/actions/workflows"
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
		workflow.Repository = repository.Url
		workflows = append(workflows, workflow)
	}

	return workflows
}
