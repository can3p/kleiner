package version

import (
	"fmt"

	"github.com/can3p/kleiner/shared/published"
	"github.com/can3p/kleiner/shared/types"
	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New(buildinfo *types.BuildInfo) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(ccmd *cobra.Command, args []string) error {
			return run(buildinfo)
		},
	}

	return versionCmd
}

func run(buildinfo *types.BuildInfo) (err error) {
	_, err = fmt.Println(buildinfo)

	if err != nil {
		return
	}

	upstreamVersion, err := published.GetLastPublishedVersion(buildinfo)

	if err != nil {
		return err
	}

	fmt.Printf("Latest available version: %s\n", upstreamVersion)

	return
}
