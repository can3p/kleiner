//go:build production

package buildinfo

import (
	"time"

	"github.com/can3p/kleiner/shared/version"
)

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
