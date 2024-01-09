package published

import (
	"fmt"

	"github.com/can3p/kleiner/generated/buildinfo"
	"github.com/can3p/kleiner/generated/internal/version"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v57/github"
)

func getAPIRelaseUrl() string {
	repo := buildinfo.GithubRepo()

	return fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
}

func GetLastPublishedVersion() (*version.Version, error) {
	var releaseObject github.RepositoryRelease

	client := resty.New()

	_, err := client.R().SetResult(&releaseObject).Get(getAPIRelaseUrl())

	if err != nil {
		return nil, err
	}

	parsed, err := version.Parse(*releaseObject.TagName)

	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
