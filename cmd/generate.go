package cmd

import (
	"log"
	"os"

	"github.com/can3p/go-scarf/scaffolder"
	"github.com/can3p/kleiner/internal/project"
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

	p, err := project.ResolveProjectFromCwd()

	var projectName string
	var githubRepo string

	if err != nil {
		log.Println(err)
	} else {
		projectName = p.Name
		githubRepo = p.GithubRepo
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	generateCmd.Flags().String("project-name", projectName, "Project and binary name, current folder name by default")
	generateCmd.Flags().String("github-repo", githubRepo, "Github repo in form user/repo")
	generateCmd.Flags().String("out", path, "Output folder")
	generateCmd.Flags().BoolVarP(&test, "test", "", false, "Do not write anything, write everything to stdout")

	_ = generateCmd.MarkFlagDirname("out")

	return generateCmd
} // generateCmd represents the generate command

func init() {
	rootCmd.AddCommand(GenerateCommand())
}
