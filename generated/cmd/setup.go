package cmd

import (
	"github.com/can3p/kleiner/generated/cmd/version"
	"github.com/spf13/cobra"
)

func Setup(root *cobra.Command) {
	root.AddCommand(version.New())
}
