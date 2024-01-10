package cmd

import (
	"github.com/can3p/kleiner/shared/cmd/cobra/update"
	"github.com/can3p/kleiner/shared/cmd/cobra/version"
	"github.com/can3p/kleiner/shared/types"
	"github.com/spf13/cobra"
)

func Setup(buildinfo *types.BuildInfo, root *cobra.Command) {
	root.AddCommand(version.New(buildinfo))
	root.AddCommand(update.New(buildinfo))
}
