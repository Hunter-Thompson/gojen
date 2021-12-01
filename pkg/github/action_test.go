package github_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"

	"github.com/Hunter-Thompson/gojen/pkg/github"
	"github.com/Hunter-Thompson/gojen/pkg/project"
)

func TestA(t *testing.T) {
	a := github.CreateWorkflow("asd")
	a.AddJobs(map[string]*github.Job{
		"asd": {
			Name:   project.String("asd"),
			RunsOn: project.String("ubuntu-latest"),
			Steps: &[]*github.JobStep{
				{
					Name: project.String("asdjob"),
					Run:  project.String("echo asd"),
				},
			}},
		"dsa": {
			Name:   project.String("dsa"),
			RunsOn: project.String("ubuntu-latest"),
			Steps: &[]*github.JobStep{
				{
					Name: project.String("dsajob"),
					Run:  project.String("echo asd"),
				},
			}},
	})

	err := a.AppendStep("asd", &github.JobStep{
		Name: project.String("appendtest"),
		Run:  project.String("echo asd"),
	})

	if err != nil {
		t.Error(err)
	}
	err = a.PrependStep("dsa", &github.JobStep{
		Name: project.String("prependtest"),
		Run:  project.String("echo asd"),
	})
	if err != nil {
		t.Error(err)
	}

	b, err := a.ConvertToYAML()
	if err != nil {
		t.Error(err)
	}

	err = cupaloy.Snapshot(b)
	if err != nil {
		t.Error(err)
	}

}
