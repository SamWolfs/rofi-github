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

	"github.com/cli/go-gh"
)

type Repository struct {
	Name  string `mapstructure:"name"`
	Owner struct {
		Login string `mapstructure:"login"`
	} `mapstructure:"owner"`
	Pretty string `mapstructure:"pretty"`
}

func (r Repository) Format(color string) string {
	info, _ := json.Marshal(r)
	return fmt.Sprintf("%s <span color=\"%s\" size=\"10pt\" style=\"italic\">(%s)</span>\x00info\x1f%s", r.Print(), color, r.Url(), string(info[:]))
}

func (r Repository) Print() string {
	if r.Pretty == "" {
		return r.Name
	}
	return r.Pretty
}

func (r Repository) View() {
	_, stdErr, err := gh.Exec("repo", "view", "--web", r.Url())
	if err != nil {
		log.Fatal(string(stdErr.Bytes()[:]))
	}
}

func (repository Repository) Url() string {
	return fmt.Sprintf("%s/%s", repository.Owner.Login, repository.Name)
}

func GetUserRepositories() []Repository {
	stdOut, _, err := gh.Exec("repo", "list", "--json", "name,owner")
	if err != nil {
		panic(err)
	}

	var repositories []Repository
	if err := json.Unmarshal(stdOut.Bytes(), &repositories); err != nil {
		panic(err)
	}

	return repositories
}
