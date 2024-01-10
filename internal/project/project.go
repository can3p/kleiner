package project

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Project struct {
	Name       string
	GithubRepo string
}

func ResolveProjectFromCwd() (*Project, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to resolve current working folder")
	}

	cwdFolder := filepath.Base(path)

	currentRepo, err := resolveGitRepo()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to resolve github repo")
	}

	return &Project{
		Name:       cwdFolder,
		GithubRepo: currentRepo,
	}, nil
}

func resolveGitRepo() (string, error) {
	cmd := exec.Command("git", "ls-remote", "--get-url")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	repo := out.String()
	repo = strings.TrimSpace(repo)
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")

	return repo, nil
}
