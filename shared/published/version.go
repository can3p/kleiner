package published

import (
	"fmt"

	"github.com/can3p/kleiner/shared/types"
	"github.com/can3p/kleiner/shared/version"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v57/github"
)

func getAPIRelaseUrl(buildinfo *types.BuildInfo) string {
	repo := buildinfo.GithubRepo

	return fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
}

func GetLastPublishedVersion(buildinfo *types.BuildInfo) (*version.Version, error) {
	var releaseObject github.RepositoryRelease

	client := resty.New()

	_, err := client.R().SetResult(&releaseObject).Get(getAPIRelaseUrl(buildinfo))

	if err != nil {
		return nil, err
	}

	parsed, err := version.Parse(*releaseObject.TagName)

	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
