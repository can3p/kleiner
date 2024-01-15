package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/can3p/kleiner/internal/cmdtools"
	"github.com/can3p/kleiner/internal/project"
	"github.com/can3p/kleiner/shared/published"
	"github.com/can3p/kleiner/shared/types"
	"github.com/can3p/kleiner/shared/version"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/pkg/errors"
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
			p, err := project.ResolveProjectFromCwd()

			if err != nil {
				return errors.Wrapf(err, "failed to resolve project")
			}

			var newVersion version.Version
			var lastVersionStr string

			lastVersion, err := published.GetLastPublishedVersion(p.GithubRepo)

			if err == published.ErrNoReleaseFound {
				newVersion = version.Version{
					Major: 0,
					Minor: 0,
					Patch: 1,
				}
				lastVersionStr = "not found"
			} else if err != nil {
				return err
			} else {
				newVersion = lastVersion.Increment(time.Now())
				lastVersionStr = lastVersion.String()
			}

			input := confirmation.New(
				fmt.Sprintf("Do you want to release version [%s]? Last version is %s.", newVersion, lastVersionStr),
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

			releaseSHA, err := releaseTagExists(vversion)

			if err != nil {
				return errors.Wrapf(err, "failed to check the tag")
			}

			currentSHA, err := currentSHA()

			if err != nil {
				return errors.Wrapf(err, "failed to check current SHA")
			}

			if releaseSHA != "" && currentSHA != releaseSHA {
				return errors.Errorf("Tag %s already exists and points to [%s], but the HEAD point to [%s]",
					vversion, releaseSHA, currentSHA)
			}

			if releaseSHA == "" {
				cc := fmt.Sprintf("git tag -a %s -m '%s'", vversion, tagComment)
				if err := cmdtools.RunCmd(cc); err != nil {
					fmt.Printf("Error creating a tag [%s]: %v\n", vversion, err)

					os.Exit(1)
				}
			} else {
				fmt.Printf("Tag [%s] already exists locally\n", vversion)
			}

			if err := pushTag(vversion); err != nil {
				if err == ErrReleaseExistsOnRemote {
					fmt.Printf("Tag [%s] already exists on remote\n", vversion)
				} else {
					fmt.Printf("Error pushing the tag to remote [%s]: %v\n", vversion, err)

					os.Exit(1)
				}
			}

			if err := cmdtools.RunCmd("goreleaser release --clean"); err != nil {
				fmt.Printf("Error publishing the version [%s]: %v\n", vversion, err)

				os.Exit(1)
			}

			return nil
		},
	}

	releaseCmd.Flags().StringVarP(&tagComment, "tag-comment", "m", "", "a message for the new tag")
	_ = releaseCmd.MarkFlagRequired("tag-comment")

	return releaseCmd
} // releaseCmd represents the release command

func releaseTagExists(version string) (string, error) {
	cc := fmt.Sprintf("git rev-list -1 %s", version)
	output, err := cmdtools.RunCmdAndGetOutput(cc)

	if err != nil {
		return "", err
	} else if output.ErrorCode == 128 {
		// git told us that no information has been found
		return "", nil
	} else if output.ErrorCode > 0 {
		return "", errors.Errorf("Failed to run [%s], the error code: %d, stderr: %s", cc, output.ErrorCode, string(output.Stderr))
	}

	sha := string(output.Stdout)
	sha = strings.TrimSpace(sha)

	return sha, nil
}

var ErrReleaseExistsOnRemote = errors.Errorf("Release already exists on remote")

func pushTag(version string) error {
	cc := fmt.Sprintf("git push origin %s", version)
	output, err := cmdtools.RunCmdAndGetOutput(cc)

	if err != nil {
		return err
	} else if output.ErrorCode == 0 {
		return nil
	}

	// XXX: this is fragile, needs to be replaced with a more robust check
	if strings.Contains(string(output.Stderr), "already exists in the remote") {
		return ErrReleaseExistsOnRemote
	}

	return errors.Errorf("Failed to run [%s], the error code: %d, stderr: %s", cc, output.ErrorCode, string(output.Stderr))
}

func currentSHA() (string, error) {
	output, err := cmdtools.RunCmdAndGetOutputNoErrCode("git rev-list -1 HEAD")

	if err != nil {
		return "", err
	}

	sha := string(output.Stdout)
	sha = strings.TrimSpace(sha)

	return sha, nil
}
