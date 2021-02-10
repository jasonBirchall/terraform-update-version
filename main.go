/*
Basic premise for this quick tool:
- ingest a list of repositories
- pull repository
- create branch
- run terraform upgrade tool
- commit changes
- create pr
- display link to pr at the end
*/
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	// . "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
)

const file = "example"

var (
	commitMessage = "Run Terraform 0.13 upgrade tool on repository"
	commitBranch  = "terraform-0.13-upgrade"
	sourceOwner   = "ministryofjustice"
	baseBranch    = "main"
	prSubject     = "Upgrade Terraform to 0.13.x"
	prDescription = "This PR contains the relevant changes required to upgrade the repository Terraform code to 0.13.x"
	authorName    = "jasonbirchall"
	authorEmail   = "jason.birchall@digital.justice.gov.uk"
)

var client *github.Client
var ctx = context.Background()
var url = "https://github.com/ministryofjustice/"

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorised: No token present")
	}

	repos, err := getRepos()
	if err != nil {
		log.Fatalf("Unable to find file: %s\n", err)
	}

	for _, repo := range repos {
		r, err := git.PlainClone(repo, true, &git.CloneOptions{
			// The intended use of a GitHub personal access token is in replace of your password
			// because access tokens can easily be revoked.
			// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
			Auth: &http.BasicAuth{
				Username: "jasonBirchall", // yes, this can be anything except an empty string
				Password: token,
			},
			URL:      url + repo,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatalf("Unable to clone repo: %s\n", err)
		}
		fmt.Println(r)

	}

	// 	ref, err := r.Head()
	// 	if err != nil {
	// 		log.Fatalf("Unable to retrieve branch: %s\n", err)
	// 	}

	// 	commit, err := r.CommitObject(ref.Hash())
	// 	if err != nil {
	// 		log.Fatalf("Unable to retrieve commit object: %s\n", err)
	// 	}
	// 	fmt.Println(commit)

	// 	// run command
	// 	// commit
	// 	// pull request
	// 	// add link to pr to collection and print
	// }

}

func getRepos() ([]string, error) {
	var s []string
	file, err := os.Open(file)
	if err != nil {
		return s, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return s, err

	}
	return s, nil
}
