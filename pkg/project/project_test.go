package project_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/Hunter-Thompson/gojen/pkg/github"
	"github.com/Hunter-Thompson/gojen/pkg/project"
	"github.com/bradleyjkemp/cupaloy/v2"
)

func TestSProject(t *testing.T) {
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
			WorkflowEnv: &map[string]*string{
				"asd": project.String("testenv"),
			},
			AppendSteps: &[]*github.JobStep{
				{
					Name: project.String("appendtest1"),
					Run:  project.String("test1"),
				},
				{
					Name: project.String("appendtest2"),
					Run:  project.String("test2"),
				},
			},
			PrependSteps: &[]*github.JobStep{
				{
					Name: project.String("prependteststep1"),
					Run:  project.String("test1"),
				},
				{
					Name: project.String("prependteststep2"),
					Run:  project.String("test2"),
				},
			},
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
			GoTestArgs:           project.StringSlice([]string{"-v", "-covera"}),
			CodeCov:              project.Bool(true),
			GoBuild:              project.Bool(true),
			GoBuildArgs:          project.StringSlice([]string{""}),
			WorkflowEnv: &map[string]*string{
				"asd": project.String("testenv1"),
			},
			AppendSteps: &[]*github.JobStep{
				{
					Name: project.String("appendtest2"),
					Run:  project.String("test2"),
				},
				{
					Name: project.String("appendtest3"),
					Run:  project.String("test3"),
				},
			},
			PrependSteps: &[]*github.JobStep{
				{
					Name: project.String("prependteststep4"),
					Run:  project.String("test1"),
				},
				{
					Name: project.String("prependteststep4"),
					Run:  project.String("test2"),
				},
			},
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
			WorkflowEnv: &map[string]*string{
				"asd": project.String("testenv2"),
			},
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
			WorkflowEnv: &map[string]*string{
				"asd": project.String("testenv3"),
			},
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

	for k, p := range projects {
		t.Run(strconv.Itoa(k), func(t *testing.T) {

			fmt.Println(k)

			project.CI = true

			dir, err := ioutil.TempDir("/tmp", p.GetName())
			if err != nil {
				t.Error(err.Error())
			}

			pwd, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			fmt.Println(pwd)

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

				for k, v := range *createdProject.GetWorkflowEnv() {
					if *v != *(*p.WorkflowEnv)[k] {
						t.Errorf("expected %s, got %s", *(*p.WorkflowEnv)[k], *v)
					}
				}
			}

			if !reflect.DeepEqual(createdProject.GetGitignore(), *p.Gitignore) {
				t.Errorf("expected %s, got %s", *p.Gitignore, createdProject.GetGitignore())
			}

			if !reflect.DeepEqual(createdProject.GetCodeOwners(), *p.CodeOwners) {
				t.Errorf("expected %s, got %s", *p.CodeOwners, createdProject.GetCodeOwners())
			}

			if k == 1 {
				p := fmt.Sprintf("%s/go.mod", dir)
				err := ioutil.WriteFile(p, []byte{}, 0644)

				if err != nil {
					t.Error(err)
				}
			}

			err = createdProject.SetupProject()
			if err != nil {
				if k == 1 {
					if err.Error() == "logged to stderr" {
						err = os.Chdir(pwd)
						if err != nil {
							t.Error(err.Error())
						}
						t.Skip()
					}
				}
				t.Error(err)
			}

			err = os.Chdir(pwd)
			if err != nil {
				t.Error(err.Error())
			}

			codeOwnersFile := filepath.Join(dir, ".github", "CODEOWNERS")
			if _, err := os.Stat(codeOwnersFile); os.IsNotExist(err) {
				t.Errorf("expected %s to exist", codeOwnersFile)
			}
			codeOwnersContents, err := ioutil.ReadFile(codeOwnersFile)
			if err != nil {
				t.Error(err)
			}
			err = cupaloy.SnapshotMulti(strconv.Itoa(k)+"codeowners", (codeOwnersContents))
			if err != nil {
				t.Error(err)
			}

			gitIgnoreFile := filepath.Join(dir, ".gitignore")
			if _, err := os.Stat(gitIgnoreFile); os.IsNotExist(err) {
				t.Errorf("expected %s to exist", gitIgnoreFile)
			}
			gitIgnoreContents, err := ioutil.ReadFile(gitIgnoreFile)
			if err != nil {
				t.Error(err)
			}
			err = cupaloy.SnapshotMulti(strconv.Itoa(k)+"gitignore", (gitIgnoreContents))
			if err != nil {
				t.Error(err)
			}

			if createdProject.GetName() != "test1" {
				buildWorkflowFile := filepath.Join(dir, ".github", "workflows", "build.yml")
				if _, err := os.Stat(buildWorkflowFile); os.IsNotExist(err) {
					t.Errorf("expected %s to exist", buildWorkflowFile)
				}
				buildWorkflowContents, err := ioutil.ReadFile(buildWorkflowFile)
				if err != nil {
					t.Error(err)
				}
				err = cupaloy.SnapshotMulti(strconv.Itoa(k)+"buildworkflow", (buildWorkflowContents))
				if err != nil {
					t.Error(err)
				}

				releaseWorfklowFile := filepath.Join(dir, ".github", "workflows", "release.yml")
				if _, err := os.Stat(releaseWorfklowFile); os.IsNotExist(err) {
					t.Errorf("expected %s to exist", releaseWorfklowFile)
				}
				releaseWorfklowContents, err := ioutil.ReadFile(releaseWorfklowFile)
				if err != nil {
					t.Error(err)
				}
				err = cupaloy.SnapshotMulti(strconv.Itoa(k)+"releaseworkflow", (releaseWorfklowContents))
				if err != nil {
					t.Error(err)
				}
			}

		})
	}

}

func TestFProject(t *testing.T) {
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

	for k, p := range failedProjects {
		t.Run(strconv.Itoa(k), func(t *testing.T) {
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

		})
	}
}
