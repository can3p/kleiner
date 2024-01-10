package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/can3p/go-scarf/scaffolder"
	kleinerTemplate "github.com/can3p/kleiner/template"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func GenerateCommand() *cobra.Command {
	var test bool

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate code for the new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := cmd.Flag("project-name").Value.String()
			githubRepo := cmd.Flag("github-repo").Value.String()
			out := cmd.Flag("out").Value.String()

			if githubRepo == "" {
				return errors.Errorf("Github repo is missing")
			}

			s := scaffolder.New()

			if !test {
				s = s.WithProcessor(scaffolder.FSProcessor(out))
			}

			return s.Scaffold(kleinerTemplate.Template, scaffolder.ScaffoldData{
				"ProjectName": projectName,
				"GithubRepo":  githubRepo,
			})
		},
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	cwdFolder := filepath.Base(path)

	currentRepo, _ := resolveGitRepo()

	generateCmd.Flags().String("project-name", cwdFolder, "Project and binary name, current folder name by default")
	generateCmd.Flags().String("github-repo", currentRepo, "Github repo in form user/repo")
	generateCmd.Flags().String("out", path, "Output folder")
	generateCmd.Flags().BoolVarP(&test, "test", "", false, "Do not write anything, write everything to stdout")

	_ = generateCmd.MarkFlagDirname("out")

	return generateCmd
} // generateCmd represents the generate command

func init() {
	rootCmd.AddCommand(GenerateCommand())
}

func resolveGitRepo() (string, error) {
	cmd := exec.Command("git", "ls-remote", "--get-url")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	repo := out.String()
	repo = strings.TrimSpace(repo)
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")

	return repo, nil
}
