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
	"strconv"

	"github.com/SamWolfs/rofi-github/github"
	"github.com/spf13/viper"
)

type Metadata struct {
	Repositories []Repository
}

type Repository struct {
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
}

type Workflow struct {
	Id         string `mapstructure:"id"`
	Name       string `mapstructure:"name"`
	Repository string `mapstructure:"repository"`
}

func ReadRepositories(repositories []github.Repository) []Repository {
	repos := make([]Repository, len(repositories))
	for i, repo := range repositories {
		r := fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)
		repos[i] = Repository{
			Name: r,
			Url:  r,
		}
	}
	return repos
}

func ReadWorkflows(workflows []github.Workflow) []Workflow {
	wfs := make([]Workflow, len(workflows))
	for i, workflow := range workflows {
		wfs[i] = Workflow{
			Id:         strconv.Itoa(workflow.Id),
			Name:       workflow.Name,
			Repository: workflow.Repository,
		}
	}
	return wfs
}

func formatRow(row interface{}) string {
	color := viper.GetString("fgColor")
	if color == "" {
		color = "#928374"
	}
	switch r := row.(type) {
	case Repository:
		return fmt.Sprintf("%s <span color=\"%s\" size=\"10pt\" style=\"italic\">(%s)</span>\x00info\x1f%s", r.Name, color, r.Url, r.Url)
	case Workflow:
		info, _ := json.Marshal(r)
		return fmt.Sprintf("%s <span color=\"%s\" size=\"10pt\" style=\"italic\">(%s)</span>\x00info\x1f%s", r.Name, color, r.Repository, string(info[:]))
	}
	return "Type not implemented."
}
