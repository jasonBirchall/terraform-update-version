package git

import (
	"fmt"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing"
	"github.com/go-git/go-git/plumbing/transport/http"
)

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
