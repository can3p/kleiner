package update

import (
	"fmt"

	"github.com/can3p/kleiner/shared/published"
	"github.com/can3p/kleiner/shared/types"
	"github.com/can3p/kleiner/shared/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New(buildinfo *types.BuildInfo) *cobra.Command {
	var versionOverride string

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update to the last update",
		RunE: func(ccmd *cobra.Command, args []string) error {
			return run(buildinfo, versionOverride)
		},
	}

	updateCmd.Flags().StringVarP(&versionOverride, "override-version", "o", "", "Update to a specific version instead of last one")

	return updateCmd
}

func run(buildinfo *types.BuildInfo, versionOverride string) (err error) {
	upstreamVersion, err := published.GetLastPublishedVersion(buildinfo.GithubRepo)

	if err != nil {
		return err
	}

	fmt.Printf("Latest available version: %s\n", upstreamVersion)

	if buildinfo.IsDev() && versionOverride == "" {
		fmt.Printf("With a dev version running no update will happen unless override-version flag is specified\n")
		return
	}

	versionToUpdateTo := upstreamVersion

	if versionOverride != "" {
		v, err := version.Parse(versionOverride)

		if err != nil {
			return errors.Wrapf(err, "Failed to parse version: [%s]", versionOverride)
		}

		versionToUpdateTo = &v
	}

	currentVersion := buildinfo.Version

	if versionToUpdateTo.Equal(currentVersion) {
		fmt.Printf("The cli is already on the version [%s], no changes will be made", versionToUpdateTo)
		return
	}

	fmt.Printf("Will attempt to update to version %s\n", versionToUpdateTo)

	err = published.DownloadAndReplaceBinary(buildinfo, versionToUpdateTo.String())

	if err == nil {
		fmt.Println("Update successful!")
	}

	return
}
