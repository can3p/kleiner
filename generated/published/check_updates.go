package published

import (
	"fmt"
	"os"

	"github.com/can3p/kleiner/generated/buildinfo"
	"github.com/fatih/color"
)

func MaybeNotifyAboutNewVersion() {
	// we never notify about new version in case of dev builds
	if buildinfo.IsDev() {
		return
	}

	version := buildinfo.Version()
	upstream, err := GetLastPublishedVersion()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for new version: %s", err.Error())
	}

	if upstream.Newer(version) {
		c := color.New(color.FgCyan)
		c.Fprintf(os.Stderr, "New version [%s] is available, your current version is %s. Run update command to update\n\n", upstream, version)
	}
}
