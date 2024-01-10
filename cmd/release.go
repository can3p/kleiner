package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/can3p/kleiner/shared/published"
	"github.com/can3p/kleiner/shared/types"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/spf13/cobra"
)

func ReleaseCommand(buildinfo *types.BuildInfo) *cobra.Command {
	var tagComment string

	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Create a new release",
		Long: `Create a new release and upload it to github.

		Kleiner will select next version and then run:

		* git tag -a v<next version>
		* git push origin  v<next version>
		* goreleaser release --clean`,

		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := published.GetLastPublishedVersion(buildinfo)

			if err != nil {
				return err
			}

			newVersion := version.Increment(time.Now())

			input := confirmation.New(
				fmt.Sprintf("Do you want to releaser version [%s]? Last version is %s.", newVersion, version),
				confirmation.No)

			ready, err := input.RunPrompt()
			if err != nil {
				fmt.Printf("Error: %v\n", err)

				os.Exit(1)
			}

			if !ready {
				return nil
			}

			vversion := "v" + newVersion.String()

			cmdChain := [][]string{
				{"git", "tag", "-a", vversion, "-m", tagComment},
				{"git", "push", "origin", vversion},
				{"goreleaser", "release", "--clean"},
			}

			for _, args := range cmdChain {
				if err := runCmd(args[0], args[1:]...); err != nil {
					fmt.Printf("Error running [%s]: %v\n", strings.Join(args, " "), err)

					os.Exit(1)
				}
			}

			return nil
		},
	}

	releaseCmd.Flags().StringVar(&tagComment, "tag-comment", "", "a message for the new tag")
	_ = releaseCmd.MarkFlagRequired("tag-comment")

	return releaseCmd
} // releaseCmd represents the release command

func runCmd(binary string, args ...string) error {
	cmd := exec.Command(binary, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
