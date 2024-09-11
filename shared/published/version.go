package published

import (
	"fmt"
	"net/http"

	"github.com/can3p/kleiner/shared/version"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v57/github"
	"github.com/pkg/errors"
)

var ErrNoReleaseFound = errors.Errorf("No release found")

func getAPIRelaseUrl(repo string) string {
	return fmt.Sprintf("https://api.github.com/repos%s/releases/latest", repo)
}

func GetLastPublishedVersion(githubRepo string) (*version.Version, error) {
	var releaseObject github.RepositoryRelease

	client := resty.New()

	resp, err := client.R().SetResult(&releaseObject).Get(getAPIRelaseUrl(githubRepo))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNoReleaseFound
	}

	parsed, err := version.Parse(*releaseObject.TagName)

	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
