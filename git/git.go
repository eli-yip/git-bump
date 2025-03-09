package git

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Songmu/gitconfig"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct{ repo *git.Repository }

func (r *Repository) GetTags() ([]string, error) {
	tagrefs, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	var tags []string
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		tag := t.Name()
		if tag.IsTag() {
			tags = append(tags, tag.Short())
		}
		return nil
	})
	return tags, err
}

func FindGitRoot(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			return path, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", errors.New("fatal: not a git repository (or any parent up to root)")
		}
		path = parent
	}
}

func Open(path string) (*Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return nil, errors.New("fatal: not a git repository")
		}
		return nil, err
	}
	return &Repository{repo: r}, nil
}

func (r *Repository) CreateTag(tag string) error {
	head, err := r.repo.Head()
	if err != nil {
		return err
	}

	user, err := gitconfig.User()
	if err != nil {
		return err
	}

	email, err := gitconfig.Email()
	if err != nil {
		return err
	}

	_, err = r.repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  user,
			Email: email,
			When:  time.Now(),
		},
		Message: tag,
	})
	return err
}
