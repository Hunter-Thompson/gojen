package project_test

import (
	"bytes"
	"testing"

	"github.com/Hunter-Thompson/gojen/pkg/project"
	"github.com/bradleyjkemp/cupaloy/v2"
)

func TestLog(t *testing.T) {

	strs := []string{
		"test",
		"test2",
		"test3",
		"test4",
	}

	strs2 := []string{
		"func",
		"func2",
		"func3",
		"func4",
	}

	for _, str := range strs {
		for _, str2 := range strs2 {

			t.Run(str, func(t *testing.T) {
				a := bytes.Buffer{}
				project.LogFail(&a, str, str2)
				project.LogInfo(&a, str, str2)
				project.LogSuccess(&a, str, str2)

				err := cupaloy.SnapshotMulti(str+str2, a.String())
				t.Log(a.String())
				if err != nil {
					t.Error(err)
				}
			})
		}
	}

}
