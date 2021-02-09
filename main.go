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
	"errors"
	"fmt"
	"log"
	"os"

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

func main() {
	repos, err := getRepos()
	if err != nil {
		log.Fatalf("Unable to find file: %s\n", err)
	}

	client = github.NewClient(nil)

	for _, r := range repos {
		ref, err := getRef(r)
		if err != nil {
			log.Fatalf("Unable to get/create the commit reference: %s\n", err)
		}
		if ref == nil {
			log.Fatalf("No error where returned but the reference is nil")
		}
		fmt.Println(ref)
	}

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

func getRef(r string) (ref *github.Reference, err error) {
	if ref, _, err := client.Git.GetRef(ctx, sourceOwner, r, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	if commitBranch == baseBranch {
		return nil, errors.New("The commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if baseBranch == "" {
		return nil, errors.New("The `base-branch` should not be set to an empty string when the branch specified by `commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, sourceOwner, r, "refs/heads/"+baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, sourceOwner, r, newRef)

	return ref, err
}
