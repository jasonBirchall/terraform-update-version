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
			log.Fatalf("Issue detected: %s\n", err)
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
			Username: user,
			Password: token,
		},
		URL: url + repo,
	})
	if err != nil {
		return err
		// if err == git.ErrRepositoryAlreadyExists {
		// 	fmt.Println("repo was already cloned")
		// } else {
		// 	return err
		// }
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	branch := "refs/heads/tf-0.13upgrade"
	b := plumbing.ReferenceName(branch)

	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: b,
	})
	if err != nil {
		return err
	}

	err = executeCommand(repo)
	if err != nil {
		return err
	}

	// Add to staging
	err = w.AddWithOptions(&git.AddOptions{
		All:  true,
		Glob: ".",
	})
	if err != nil {
		return err
	}

	// git commit -m $message
	w.Commit("Added my new file", &git.CommitOptions{})
	// commit, err := w.Commit("Added my new file", &git.CommitOptions{
	// 	All: true,
	// })
	// if err != nil {
	// 	return err
	// }

	// obj, err := r.CommitObject(commit)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(obj)
	// status, _ := w.Status()

	// fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	// commit, err := w.Commit("example go-git commit", &git.CommitOptions{})

	// // Prints the current HEAD to verify that all worked well.
	// obj, _ := r.CommitObject(commit)

	// fmt.Println(obj)
	return nil
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
