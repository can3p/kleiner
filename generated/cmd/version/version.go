// Package version implements the version command chain.
package version

import (
	"fmt"

	"github.com/can3p/kleiner/generated/buildinfo"
	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Show version information"
	)

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

	return
}
