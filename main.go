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
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	user := os.Getenv("GITHUB_AUTH_USER")
	if token == "" || user == "" {
		log.Fatal("Unauthorised: No user or token present")
	}

	repos, err := getRepos()
	if err != nil {
		log.Fatalf("Unable to find file: %s\n", err)
	}

	for _, repo := range repos {
		err := cloneRepo(repo, token, user)
		if err != nil {
			log.Fatalf("Unable to clone repository: %s\n", err)
		}

		err = executeCommand(repo)
		if err != nil {
			log.Fatalf("Unable to execute command: %s\n", err)
		}
	}
}

func executeCommand(repo string) error {
	cmd := execute.ExecTask{
		Command:     "terraform",
		Args:        []string{"0.13upgrade", "--yes"},
		StreamStdio: false,
	}

	os.Chdir(repo)
	_, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = cmd.Execute()
	if err != nil {
		return err
	}

	os.Chdir("..")
	_, err = os.Getwd()
	if err != nil {
		return err
	}

	return nil
}

func cloneRepo(repo, token, user string) error {
	r, err := git.PlainClone(repo, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: user, // yes, this can be anything except an empty string
			Password: token,
		},
		URL: url + repo,
		// Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	headRef, err := r.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference("refs/heads/tf-0.13upgrade", headRef.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(ref),
	})
	if err != nil {
		return err
	}

	return nil
}

// 	// run command
// 	// commit
// 	// pull request
// 	// add link to pr to collection and print
// }

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
