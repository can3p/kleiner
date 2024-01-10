//go:build production

package buildinfo

import (
	"time"

	"github.com/can3p/kleiner/shared/version"
)

// The code has been borrowed from flyctl project - https://github.com/superfly/flyctl/blob/0dff860a878e2b280f2f53ce2aaf21ce39d800c2/internal/buildinfo
// This code in the file is subject to Apache-2.0 license as per flyctl project

var (
	buildDate    = "<date>"
	buildVersion = "<version>"
	environment  = "production"
)

func loadBuildTime() (err error) {
	cachedBuildTime, err = time.Parse(time.RFC3339, buildDate)
	return
}

func loadVersion() (err error) {
	cachedVersion, err = version.Parse(buildVersion)
	return
}
