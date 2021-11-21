/*
Copyright © 2021 Aatman <aatman@auroville.org.in>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Hunter-Thompson/gojen/pkg/project"
	"github.com/spf13/cobra"
)

var cfg project.Project

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}

		if _, err := os.Stat(pwd + "/gojen.json"); errors.Is(err, os.ErrNotExist) {
			err := cfg.WriteConfig()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			cmd.Println("config written")

			proj, err := project.InitProject()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			err = proj.SetupProject()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("initialized new project with name %s\n", cfg.GetName())
			return
		}

		fmt.Println("config already exists")

		proj, err := project.InitProject()
		if err != nil {
			fmt.Println(err.Error())
		}

		err = proj.SetupProject()
		if err != nil {
			fmt.Println(err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	cfg.Name = newCmd.Flags().String("name", "", "name of your project")
	cfg.Description = newCmd.Flags().String("description", "", "description of your project")
	cfg.Repository = newCmd.Flags().String("repository", "", "repository URL for your project")
	cfg.GoVersion = newCmd.Flags().String("goversion", "1.16", "go version")
	cfg.AuthorName = newCmd.Flags().String("authorname", "", "name of Author")
	cfg.AuthorEmail = newCmd.Flags().String("authoremail", "", "email of Author")
	cfg.AuthorOrganization = newCmd.Flags().String("authororganization", "", "author github organization")
	cfg.Licensed = newCmd.Flags().Bool("licensed", false, "add a license")
	cfg.Release = newCmd.Flags().Bool("release", false, "setup go-semantic-release and upload binaries")
	cfg.DefaultReleaseBranch = newCmd.Flags().String("defaultreleasebranch", "master", "default branch to release from")
	cfg.Gitignore = newCmd.Flags().StringSlice("gitignore", []string{}, "list of gitignore values")
	cfg.CodeOwners = newCmd.Flags().StringSlice("codeowners", []string{}, "list of codeowner values")
	cfg.GoLinter = newCmd.Flags().Bool("golinter", false, "enable golinter")
	cfg.GoTest = newCmd.Flags().Bool("gotest", true, "enable go test")
	cfg.GoTestArgs = newCmd.Flags().StringSlice("gotestargs", []string{}, "arguments for go test")
}
