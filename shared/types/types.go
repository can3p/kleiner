package types

import (
	"fmt"
	"time"

	"github.com/can3p/kleiner/shared/version"
)

// The struct has been borrowed from flyctl project - https://github.com/superfly/flyctl/blob/0dff860a878e2b280f2f53ce2aaf21ce39d800c2/internal/buildinfo/buildinfo.go
// This code in the file is subject to Apache-2.0 license as per flyctl project

type BuildInfo struct {
	Name         string
	Version      version.Version
	Commit       string
	BranchName   string
	BuildDate    time.Time
	OS           string
	Architecture string
	Environment  string
	GithubRepo   string
	ProjectName  string
}

func (i BuildInfo) String() string {
	res := fmt.Sprintf("%s v%s %s/%s Commit: %s BuildDate: %s",
		i.Name,
		i.Version,
		i.OS,
		i.Architecture,
		i.Commit,
		i.BuildDate.Format(time.RFC3339))
	if i.BranchName != "" {
		res += fmt.Sprintf(" BranchName: %s", i.BranchName)
	}

	res += fmt.Sprintf(" Github Repo: https://github.com/%s", i.GithubRepo)

	return res
}

// from https://github.com/superfly/flyctl/blob/0dff860a878e2b280f2f53ce2aaf21ce39d800c2/internal/buildinfo/env.go

func (i BuildInfo) Env() string {
	return i.Environment
}

func (i BuildInfo) IsDev() bool {
	return i.Environment == "development"
}

func (i BuildInfo) IsRelease() bool {
	return !i.IsDev()
}
