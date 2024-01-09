// Package version implements the version command chain.
package version

import (
	"fmt"

	"github.com/can3p/kleiner/generated/buildinfo"
	"github.com/can3p/kleiner/generated/published"
	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE:  run,
	}

	return versionCmd
}

func run(ccmd *cobra.Command, args []string) (err error) {
	var (
		info = buildinfo.Info()
	)

	_, err = fmt.Println(info)

	if err != nil {
		return
	}

	upstreamVersion, err := published.GetLastPublishedVersion()

	if err != nil {
		return err
	}

	fmt.Printf("Latest available version: %s\n", upstreamVersion)

	return
}
