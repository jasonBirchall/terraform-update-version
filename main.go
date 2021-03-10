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
	"path/filepath"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/go-git/go-git/plumbing"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/google/go-github/github"
	git "github.com/jasonbirchall/terraform-update-version/pkg/api"
)

const (
	file          = "repositories"
	url           = "https://github.com/ministryofjustice/"
	commitSummary = "Evaluate the repository with terraform0.13upgrade binary"
	commitBody    = "In our vision to upgrade Terraform to 0.13, we will need to run this binary over any Terraform files in our repository."
	prSummary     = "Terraform 0.13 upgrade for repository"
	prBody        = "This PR contains all changes in this repository made by the command `terraform 0.13upgrade` tool."
)

var (
	client *github.Client
	ctx    = context.Background()
)

func main() {
	var s, f []string

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
		fmt.Println("--- Starting", repo)
		err := gitCreate(repo, token, user)
		if err != nil {
			log.Printf("Issue creating git repo or branch: %s\n", err)
		}

		err = gitPush()
		if err != nil {
			f = append(f, repo)
			log.Printf("Issue pushing, pulling or PR: %s\n", err)
		} else {
			s = append(s, repo)
		}

		fmt.Println("--- Finishing", repo)

		err = remove(repo)
		if err != nil {
			log.Printf("Error deleting the dir: %s", err)
		}

	}

	fmt.Println("Sucessful repos:", s)
	fmt.Println("Failed repos:", f)
}

func gitCreate(repo, token, user string) error {
	r, err := git.PlainClone(repo, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: user,
			Password: token,
		},
		URL: url + repo,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Println("worktree")
		return err
	}

	branch := "refs/heads/tf-0.13"
	b := plumbing.ReferenceName(branch)

	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: b,
	})
	if err != nil {
		return err
	}

	err = walkExecute(repo)
	if err != nil {
		return err
	}

	err = w.AddWithOptions(&git.AddOptions{
		All:  true,
		Glob: ".",
	})
	if err != nil {
		return err
	}

	err = chDir(repo)
	if err != nil {
		return err
	}

	return nil
}

func walkMatch(start string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func tfUpgrade() error {
	cmd := execute.ExecTask{
		Command:     "terraform0.13",
		Args:        []string{"0.13upgrade", "--yes"},
		StreamStdio: false,
	}

	_, err := cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func remove(r string) error {
	cmd := execute.ExecTask{
		Command:     "rm",
		Args:        []string{"-rf", r},
		StreamStdio: false,
	}

	_, err := cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func chDir(path string) error {
	os.Chdir(path)
	_, err := os.Getwd()
	if err != nil {
		return err
	}

	return nil
}

func walkExecute(repo string) error {
	dirs, err := walkMatch(repo)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		err = chDir(dir)
		if err != nil {
			return err
		}

		err = tfUpgrade()
		if err != nil {
			return err
		}

		err = chDir("/app")
		if err != nil {
			return err
		}

	}
	return nil
}

func commit() error {
	commit := execute.ExecTask{
		Command:     "git",
		Args:        []string{"commit", "-m", commitSummary, "-m", commitBody},
		StreamStdio: false,
	}

	_, err := commit.Execute()
	if err != nil {
		return err
	}
	return nil
}

func push() error {
	push := execute.ExecTask{
		Command:     "git",
		Args:        []string{"push", "--set-upstream", "origin", "tf-0.13"},
		StreamStdio: true,
	}

	_, err := push.Execute()
	if err != nil {
		return err
	}

	return nil
}

func pullRequest() error {
	pr := execute.ExecTask{
		Command:     "gh",
		Args:        []string{"pr", "create", "-t", prSummary, "-b", prBody},
		StreamStdio: true,
	}

	_, err := pr.Execute()
	if err != nil {
		return err
	}

	return nil
}

func gitPush() error {
	err := commit()
	if err != nil {
		return err
	}

	err = push()
	if err != nil {
		return err
	}

	err = pullRequest()
	if err != nil {
		return err
	}

	err = chDir("..")
	if err != nil {
		return err
	}

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
