package project

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Project struct {
	Name       string
	GithubRepo string
	GitBranch  string
}

func ResolveProjectFromCwd() (*Project, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to resolve current working folder")
	}

	cwdFolder := filepath.Base(path)

	currentRepo, err := resolveGitRepo()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to resolve github repo, it will have to be specified manually")
	}

	currentBranch, err := resolveGitBranch()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to resolve current git branch, it will have to be specified manually")
	}

	return &Project{
		Name:       cwdFolder,
		GithubRepo: currentRepo,
		GitBranch:  currentBranch,
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
	if repo[0] == byte('/') { // we are in happy ASCII land
		repo = repo[1:] // remove leading slash if it exists
	}

	return repo, nil
}

func resolveGitBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--no-color", "--show-current")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	branch := out.String()
	branch = strings.TrimSpace(branch)

	return branch, nil
}
