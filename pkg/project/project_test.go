package project_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/Hunter-Thompson/gojen/pkg/project"
)

func TestProject(t *testing.T) {
	t.Log("test")

	// TODO: add tests for codeowners, workflow, gitignore, main.go files

	projects := []project.Project{
		{
			Name:                 project.String("test"),
			Description:          project.String("test"),
			Repository:           project.String("github.com/test/test"),
			AuthorName:           project.String("test"),
			AuthorEmail:          project.String("test"),
			AuthorOrganization:   project.String("test"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(true),
			DefaultReleaseBranch: project.String("test"),
			BuildWorkflow:        project.Bool(true),
			Gitignore:            project.StringSlice([]string{"test", "test"}),
			CodeOwners:           project.StringSlice([]string{"test", "test"}),
			GoLinter:             project.Bool(true),
			GoTest:               project.Bool(false),
			GoTestArgs:           project.StringSlice([]string{"-v", "-cover", "./..."}),
			GoBuild:              project.Bool(true),
			GoBuildArgs:          project.StringSlice([]string{""}),
		},
		{
			Name:                 project.String("test1"),
			Description:          project.String("test1"),
			Repository:           project.String("github.com/test/test1"),
			AuthorName:           project.String("test1"),
			AuthorEmail:          project.String("test1"),
			AuthorOrganization:   project.String("test1"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(false),
			DefaultReleaseBranch: project.String("test1"),
			Gitignore:            project.StringSlice([]string{"test1", "test1"}),
			CodeOwners:           project.StringSlice([]string{"test1", "test1"}),
			GoLinter:             project.Bool(false),
			GoTest:               project.Bool(true),
			GoTestArgs:           project.StringSlice([]string{"-v", "-cover"}),
			CodeCov:              project.Bool(true),
			GoBuild:              project.Bool(true),
			GoBuildArgs:          project.StringSlice([]string{""}),
		},
		{
			Name:                 project.String("test2"),
			Description:          project.String("test2"),
			Repository:           project.String("github.com/test/test2"),
			AuthorName:           project.String("test2"),
			BuildWorkflow:        project.Bool(true),
			AuthorEmail:          project.String("test2"),
			AuthorOrganization:   project.String("test2"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(true),
			DefaultReleaseBranch: project.String("test2"),
			Gitignore:            project.StringSlice([]string{"test2", "test2"}),
			CodeOwners:           project.StringSlice([]string{"test2", "test2"}),
			GoLinter:             project.Bool(true),
			GoTest:               project.Bool(true),
			GoTestArgs:           project.StringSlice([]string{"", "-cover", "./..."}),
			CodeCov:              project.Bool(true),
			GojenVersion:         project.String("1.2.0"),
			GoBuild:              project.Bool(false),
			GoBuildArgs:          project.StringSlice([]string{""}),
		},
		{
			Name:                 project.String("test3"),
			Description:          project.String("test3"),
			AuthorName:           project.String("test3"),
			Repository:           project.String("github.com/test/test3"),
			BuildWorkflow:        project.Bool(true),
			AuthorEmail:          project.String("test3"),
			AuthorOrganization:   project.String("test3"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(true),
			DefaultReleaseBranch: project.String("test3"),
			Gitignore:            project.StringSlice([]string{"test3", "test3"}),
			CodeOwners:           project.StringSlice([]string{"test3", "test3"}),
			GoLinter:             project.Bool(false),
			GoTest:               project.Bool(false),
			GoTestArgs:           project.StringSlice([]string{"-v", "-cover", "./..."}),
			GojenVersion:         project.String("1.2.0"),
			GoBuild:              project.Bool(true),
			GoBuildArgs:          project.StringSlice([]string{""}),
		},
		{
			Name:                 project.String("test4"),
			Description:          project.String("test4"),
			Repository:           project.String("github.com/test/test4"),
			AuthorName:           project.String("test4"),
			AuthorEmail:          project.String("test4"),
			AuthorOrganization:   project.String("test4"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(true),
			DefaultReleaseBranch: project.String("test4"),
			BuildWorkflow:        project.Bool(true),
			Gitignore:            project.StringSlice([]string{"test", "test"}),
			CodeOwners:           project.StringSlice([]string{"test", "test"}),
			GoLinter:             project.Bool(true),
		},
	}

	for _, p := range projects {

		project.CI = true

		dir, err := ioutil.TempDir("/tmp", p.GetName())
		if err != nil {
			t.Error(err.Error())
		}

		defer os.RemoveAll(dir)

		err = os.Chdir(dir)
		if err != nil {
			t.Error(err.Error())
		}

		err = p.WriteConfig()
		if err != nil {
			t.Error(err.Error())
		}

		createdProject, err := project.InitProject()
		if err != nil {
			t.Error(err.Error())
		}

		if createdProject.GetDescription() != *p.Description {
			t.Errorf("expected %s, got %s", *p.Description, createdProject.GetDescription())
		}

		if createdProject.GetAuthorName() != *p.AuthorName {
			t.Errorf("expected %s, got %s", *p.AuthorName, createdProject.GetAuthorName())
		}

		if createdProject.GetAuthorEmail() != *p.AuthorEmail {
			t.Errorf("expected %s, got %s", *p.AuthorEmail, createdProject.GetAuthorEmail())
		}

		if createdProject.GetAuthorOrganization() != *p.AuthorOrganization {
			t.Errorf("expected %s, got %s", *p.AuthorOrganization, createdProject.GetAuthorOrganization())
		}

		if createdProject.IsLicensed() != *p.Licensed {
			t.Errorf("expected %t, got %t", *p.Licensed, createdProject.IsLicensed())
		}

		if createdProject.IsRelease() != *p.Release {
			t.Errorf("expected %t, got %t", *p.Release, createdProject.IsRelease())
		}

		if createdProject.GetDefaultReleaseBranch() != *p.DefaultReleaseBranch {
			t.Errorf("expected %s, got %s", *p.DefaultReleaseBranch, createdProject.GetDefaultReleaseBranch())
		}

		if createdProject.IsGoLinter() != *p.GoLinter {
			t.Errorf("expected %t, got %t", *p.GoLinter, createdProject.IsGoLinter())
		}

		if *p.Name != "test4" {

			if createdProject.IsGoTest() != *p.GoTest {
				t.Errorf("expected %t, got %t", *p.GoTest, createdProject.IsGoTest())
			}

			if createdProject.IsGoBuild() != *p.GoBuild {
				t.Errorf("expected %t, got %t", *p.GoBuild, createdProject.IsGoBuild())
			}

			if !reflect.DeepEqual(createdProject.GetGoTestArgs(), *p.GoTestArgs) {
				t.Errorf("expected %s, got %s", *p.GoTestArgs, createdProject.GetGoTestArgs())
			}

			if !reflect.DeepEqual(createdProject.GetGoBuildArgs(), *p.GoBuildArgs) {
				t.Errorf("expected %s, got %s", *p.GoBuildArgs, createdProject.GetGoBuildArgs())
			}
		}

		if !reflect.DeepEqual(createdProject.GetGitignore(), *p.Gitignore) {
			t.Errorf("expected %s, got %s", *p.Gitignore, createdProject.GetGitignore())
		}

		if !reflect.DeepEqual(createdProject.GetCodeOwners(), *p.CodeOwners) {
			t.Errorf("expected %s, got %s", *p.CodeOwners, createdProject.GetCodeOwners())
		}

		err = createdProject.SetupProject()
		if err != nil {
			// a := strings.Contains(err.Error(), "package asd is not in GOROOT")
			// if !a {
			// 	t.Error(err)
			// }
			t.Error(err)
		}

	}

	failedProjects := []project.Project{
		{
			Name:                 project.String("test4"),
			Description:          project.String("test4"),
			AuthorName:           project.String("test4"),
			AuthorEmail:          project.String("test4"),
			AuthorOrganization:   project.String("test4"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(false),
			DefaultReleaseBranch: project.String("test4"),
			Gitignore:            project.StringSlice([]string{"test4", "test4"}),
			CodeOwners:           project.StringSlice([]string{"test4", "test4"}),
			GoLinter:             project.Bool(false),
			GoTest:               project.Bool(false),
			GoTestArgs:           project.StringSlice([]string{"-v", "-cover", "./..."}),
		},
		{
			Description:          project.String("test5"),
			Repository:           project.String("github.com/test/test5"),
			AuthorName:           project.String("test5"),
			AuthorEmail:          project.String("test5"),
			AuthorOrganization:   project.String("test5"),
			Licensed:             project.Bool(true),
			Release:              project.Bool(false),
			DefaultReleaseBranch: project.String("test5"),
			Gitignore:            project.StringSlice([]string{"test4", "test4"}),
			CodeOwners:           project.StringSlice([]string{"test4", "test4"}),
			GoLinter:             project.Bool(false),
			GoTest:               project.Bool(false),
			GoTestArgs:           project.StringSlice([]string{"-v", "-cover", "./..."}),
		},
	}

	for _, p := range failedProjects {
		project.CI = true
		dir, err := ioutil.TempDir("/tmp", "failedtests")
		if err != nil {
			t.Error(err.Error())
		}

		err = os.Chdir(dir)
		if err != nil {
			t.Error(err.Error())
		}

		err = p.WriteConfig()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = project.InitProject()
		if err == nil {
			fmt.Println(p.GetDescription())
			t.Error("expected error, got nil")
		}

	}

}
