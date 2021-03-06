package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/Hunter-Thompson/gojen/pkg/github"
	"github.com/Hunter-Thompson/gojen/pkg/license"
)

var CI bool

type IProject interface {
	WriteConfig() error
	SetupProject() error
	SetGitignore() error
	CreateReadme() error
	RunTest() error
	RunBuild() error
	RunLinter() error
	setCommonJobs(wf github.IAction) (github.IAction, error)
	getCommonSteps() []*github.JobStep

	GetName() string
	GetDescription() string
	GetRepository() string
	GetAuthorName() string
	GetAuthorEmail() string
	GetAuthorOrganization() string
	IsRelease() bool
	IsGoBuild() bool
	GetDefaultReleaseBranch() string
	GetGitignore() []string
	GetCodeOwners() []string
	IsGoLinter() bool
	IsGoTest() bool
	GetGoTestArgs() []string
	IsBuildWorkflow() bool
	GetGitHubToken() string
	GetGojenVersion() string
	IsIsGojen() bool
	IsCreateReadme() bool
	IsCodeCov() bool
	GetGoBuildArgs() []string
	GetWorkflowEnv() *map[string]*string
	GetLicense() string
}

type Project struct {
	Name        *string `yaml:"name" json:"name"`
	Description *string `yaml:"description" json:"description"`
	Repository  *string `yaml:"repository"  json:"repository"`
	GoVersion   *string `yaml:"goVersion" json:"goVersion"`

	AuthorName         *string `yaml:"authorName" json:"authorName"`
	AuthorEmail        *string `yaml:"authorEmail" json:"authorEmail"`
	AuthorOrganization *string `yaml:"authorOrganization" json:"authorOrganization"`

	Readme       *bool   `yaml:"readme" json:"readme"`
	GojenVersion *string `yaml:"gojenVersion" json:"gojenVersion"`
	License      *string `yaml:"license" json:"license"`

	Release              *bool     `yaml:"release" json:"release"`
	BuildWorkflow        *bool     `yaml:"buildWorkflow" json:"buildWorkflow"`
	GithubToken          *string   `yaml:"githubToken" json:"githubToken"`
	DefaultReleaseBranch *string   `yaml:"defaultReleaseBranch" json:"defaultReleaseBranch"`
	IsGojen              *bool     `yaml:"isGojen" json:"isGojen"`
	CodeCov              *bool     `yaml:"codeCov" json:"codeCov"`
	TestEnvVars          *[]string `yaml:"testEnvVars" json:"testEnvVars"`

	Gitignore  *[]string `yaml:"gitignore" json:"gitignore"`
	CodeOwners *[]string `yaml:"codeOwners" json:"codeOwners"`

	SkipVendor *bool `yaml:"skipVendor" json:"skipVendor"`
	SkipTidy   *bool `yaml:"skipTidy" json:"skipTidy"`

	GoLinter     *bool               `yaml:"goLinter" json:"goLinter"`
	GoTest       *bool               `yaml:"goTest" json:"goTest"`
	GoTestArgs   *[]string           `yaml:"goTestArgs" json:"goTestArgs"`
	GoBuild      *bool               `yaml:"goBuild" json:"goBuild"`
	GoBuildArgs  *[]string           `yaml:"goBuildArgs" json:"goBuildArgs"`
	WorkflowEnv  *map[string]*string `yaml:"workflowEnv" json:"workflowEnv"`
	PrependSteps *[]*github.JobStep  `yaml:"prependSteps" json:"prependSteps"`
	AppendSteps  *[]*github.JobStep  `yaml:"apendSteps" json:"apendSteps"`
}

func InitProject() (IProject, error) {
	proj, err := GetConfig()
	if err != nil {
		return nil, err
	}

	err = proj.ValidateConfig()
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (proj *Project) ValidateConfig() error {
	if proj.GetName() == "" {
		return errors.New("name is missing in config")
	}

	if proj.GetRepository() == "" {
		return errors.New("repository is missing in config")
	}

	return nil
}

func GetConfig() (*Project, error) {
	proj := &Project{}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfgPath := fmt.Sprintf("%s/gojen.json", pwd)

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (proj *Project) WriteConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgPath := fmt.Sprintf("%s/gojen.json", pwd)

	b, err := json.MarshalIndent(proj, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cfgPath, b, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) SetupProject() error {
	if proj.IsCodeCov() {
		err := proj.AddCodeCov()
		if err != nil {
			return err
		}
	}

	err := proj.AddLicense()
	if err != nil {
		return err
	}

	if proj.IsCreateReadme() {
		err := proj.CreateReadme()
		if err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(proj.CodeOwners, &[]string{}) {
		err := proj.SetCodeOwners()
		if err != nil {
			return err
		}
	}

	if proj.IsRelease() {
		err := proj.CreateReleaseWorkflow()
		if err != nil {
			return err
		}
	}

	if proj.IsBuildWorkflow() {
		err := proj.CreateBuildWorkflow()
		if err != nil {
			return err
		}
	}

	err = proj.SetGitignore()
	if err != nil {
		return err
	}

	modInit := exec.Command("go", "mod", "init", proj.GetRepository())
	vendor := exec.Command("go", "mod", "vendor")
	tidy := exec.Command("go", "mod", "tidy")
	gofmt := exec.Command("go", "fmt")

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(pwd + "/main.go"); errors.Is(err, os.ErrNotExist) {
		c := `package main
import (
	"fmt"
)

func main () {
	fmt.Println("project created with gojen, have fun :-)")
}`
		err := ioutil.WriteFile(pwd+"/main.go", []byte(c+"\n"), 0o644)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if _, err := os.Stat(pwd + "/go.mod"); errors.Is(err, os.ErrNotExist) {
		LogInfo(os.Stdout, "running go mod init", "Setup")

		modInit.Stdout = os.Stdout
		modInit.Stderr = os.Stderr

		err := modInit.Run()
		if err != nil {
			LogFail(os.Stderr, "running go mod vendor init failed", "Setup")
			return errors.New("logged to stderr")
		}
	}

	if err != nil {
		return err
	}

	if proj.SkipVendor == nil || !*proj.SkipVendor {
		LogInfo(os.Stdout, "running go mod vendor", "Setup")

		vendor.Stdout = os.Stdout
		vendor.Stderr = os.Stderr

		err = vendor.Run()
		if err != nil {
			LogFail(os.Stderr, "running go mod vendor failed", "Setup")
			return errors.New("logged to stderr")
		}
	}

	if proj.SkipTidy == nil || !*proj.SkipTidy {
		LogInfo(os.Stdout, "running go mod tidy", "Setup")

		tidy.Stdout = os.Stdout
		tidy.Stderr = os.Stderr

		err = tidy.Run()
		if err != nil {
			LogFail(os.Stderr, "running go mod tidy failed", "Setup")
			return errors.New("logged to stderr")
		}
	}

	LogInfo(os.Stdout, "running go fmt", "Setup")

	gofmt.Stdout = os.Stdout
	gofmt.Stderr = os.Stderr

	err = gofmt.Run()
	if err != nil {
		LogFail(os.Stderr, "running go fmt failed", "Setup")
		return errors.New("logged to stderr")
	}

	if !CI {
		if proj.IsGoLinter() {
			err = proj.RunLinter()
			if err != nil {
				return err
			}
		}
	}

	if proj.IsGoTest() {
		err = proj.RunTest()
		if err != nil {
			return err
		}
	}

	if proj.IsGoBuild() {
		err = proj.RunBuild()
		if err != nil {
			return err
		}
	}

	return nil
}

func (proj *Project) RunTest() error {
	if proj.GoTestArgs == nil {
		proj.GoTestArgs = &[]string{}
	}
	*proj.GoTestArgs = append([]string{"test"}, *proj.GoTestArgs...)

	LogInfo(os.Stdout, "running go test", "Test")

	test := exec.Command("go", *proj.GoTestArgs...)

	test.Stdout = os.Stdout
	test.Stderr = os.Stderr

	err := test.Run()
	if err != nil {
		LogFail(os.Stderr, "running go test failed", "Test")
		return errors.New("logged to stderr")
	}

	LogSuccess(os.Stdout, "go test passed", "Test")
	return nil
}

func (proj *Project) RunBuild() error {
	if proj.GoBuildArgs == nil {
		proj.GoBuildArgs = &[]string{}
	}
	*proj.GoBuildArgs = append([]string{"build"}, *proj.GoBuildArgs...)

	LogInfo(os.Stdout, "running go build", "Build")

	build := exec.Command("go", *proj.GoBuildArgs...)

	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	err := build.Run()
	if err != nil {
		LogFail(os.Stderr, "running go build failed", "Build")
		return errors.New("logged to stderr")
	}

	LogSuccess(os.Stdout, "go build passed", "Build")

	return nil
}

func (proj *Project) AddLicense() error {
	if proj.GetLicense() == "" {
		return nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(pwd + "/LICENSE"); errors.Is(err, os.ErrNotExist) {
		LogInfo(os.Stdout, "adding license", "Setup")
		licenses := []string{
			"Apache-2.0",
			"Artistic-2.0",
			"GPL-3.0-WITH-GCC-exception-3.1",
			"LGPL-2.1-or-later  ",
			"MIT-0",
			"MPL-2.0",
			"PHP-3.01",
			"Unlicense",
			"ZPL-2.1",
			"Artistic-1.0",
			"GPL-2.0-or-later",
			"GPL-3.0-or-later",
			"LGPL-3.0-or-later",
			"MIT",
			"OFL-1.1",
			"Ruby",
			"WTFPL",
		}

		if Contains(licenses, proj.GetLicense()) {
			b, err := license.Asset("license-text/" + proj.GetLicense() + ".txt")
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(pwd+"/LICENSE", b, 0o644)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (proj *Project) SetGitignore() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitignorePath := fmt.Sprintf("%s/.gitignore", pwd)
	*proj.Gitignore = append(*proj.Gitignore, *proj.Name)
	contents := strings.Join(*proj.Gitignore, "\n")

	err = ioutil.WriteFile(gitignorePath, []byte(contents), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) AddCodeCov() error {
	*proj.Gitignore = append(*proj.Gitignore, "coverage.txt")
	*proj.GoTestArgs = append([]string{"-coverprofile=coverage.txt", "-covermode=atomic"}, *proj.GoTestArgs...)

	return nil
}

func (proj *Project) CreateReadme() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	readmePath := fmt.Sprintf("%s/README.md", pwd)

	if _, err := os.Stat(readmePath); errors.Is(err, os.ErrNotExist) {
		c := `# ` + *proj.Name + `

`
		err := ioutil.WriteFile(readmePath, []byte(c+"\n"), 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proj *Project) SetCodeOwners() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/.github/CODEOWNERS", pwd)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(fmt.Sprintf("%s/.github", pwd), 0o755)
		if err != nil {
			return err
		}
	}

	contents := strings.Join(*proj.CodeOwners, "\n")

	err = ioutil.WriteFile(path, []byte(contents), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) RunLinter() error {
	LogInfo(os.Stdout, "running go linter", "Lint")

	lint := exec.Command("golangci-lint", "run")
	lint.Stdout = os.Stdout
	lint.Stderr = os.Stderr

	err := lint.Run()
	if err != nil {
		LogFail(os.Stderr, "running golint failed", "Lint")
		return errors.New("logged to stderr")
	}

	LogSuccess(os.Stdout, "go linter passed", "Lint")
	return nil
}

func (proj *Project) getCommonSteps() []*github.JobStep {
	wf := []*github.JobStep{}

	wf = append(wf, &github.JobStep{
		Name: String("Setup go"),
		Uses: String("actions/setup-go@v2"),
		With: &map[string]interface{}{
			"go-version": proj.GetGoVersion(),
		},
	})

	if proj.PrependSteps != nil {
		wf = append(wf, *proj.PrependSteps...)
	}

	if proj.IsIsGojen() {
		wf = append(wf, &github.JobStep{
			Name: String("Build and run gojen"),
			Run:  String("go build && ./gojen --ci"),
			Env:  proj.GetWorkflowEnv(),
		})
	} else {
		wf = append(wf, &github.JobStep{
			Name: String("Install gojen"),
			Run: String(fmt.Sprintf("go install github.com/Hunter-Thompson/gojen@%s",
				proj.GetGojenVersion())),
		})
		wf = append(wf, &github.JobStep{
			Name: String("Run gojen"),
			Run:  String("gojen --ci"),
			Env:  proj.GetWorkflowEnv(),
		})
	}

	if proj.IsCodeCov() {
		wf = append(wf, &github.JobStep{
			Name: String("Upload codecov coverage"),
			Uses: String("codecov/codecov-action@v2"),
			With: &map[string]interface{}{
				"files": String("./coverage.txt"),
			},
		})
	}

	wf = append(wf, &github.JobStep{
		Name: String("Check for changes"),
		Id:   String("git_diff"),
		Run:  String("git diff --exit-code || echo \"::set-output name=has_changes::true\""),
	})

	if proj.AppendSteps != nil {
		wf = append(wf, *proj.AppendSteps...)
	}

	return wf
}

func (proj *Project) setCommonJobs(wf github.IAction) (github.IAction, error) {
	wf.AddJobs(map[string]*github.Job{
		"golangci": {
			Name:   String("lint"),
			RunsOn: String("ubuntu-latest"),
			Steps: &[]*github.JobStep{
				{
					Name: String("Checkout"),
					Uses: String("actions/checkout@v2"),
				},
				{
					Name: String("Lint using golangci-lint"),
					Uses: String("golangci/golangci-lint-action@v2"),
					With: &map[string]interface{}{
						"args": String("--timeout=5m"),
					},
				},
			},
		},
	})

	wf.AddJobs(map[string]*github.Job{
		"build": {
			Name:   String("build"),
			RunsOn: String("ubuntu-latest"),
			Steps: &[]*github.JobStep{
				{
					Name: String("Checkout"),
					Uses: String("actions/checkout@v2"),
				},
			},
		},
	})

	j := proj.getCommonSteps()

	for _, v := range j {
		err := wf.AddStep("build", v)
		if err != nil {
			return nil, err
		}
	}

	return wf, nil
}

func (proj *Project) CreateReleaseWorkflow() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.github/workflows", pwd), 0o755)
	if err != nil {
		return err
	}

	wf := github.CreateWorkflow("release")

	wf.AddTrigger(github.Triggers{
		Push: &github.PushOptions{
			Branches: &[]*string{
				proj.DefaultReleaseBranch,
			},
		},
	})

	wf, err = proj.setCommonJobs(wf)
	if err != nil {
		return err
	}

	err = wf.AddStep("build", &github.JobStep{
		Name: String("Exit 1 if changes found"),
		If:   String("steps.git_diff.outputs.has_changes"),
		Run:  String("exit 1"),
	})
	if err != nil {
		return err
	}

	wf.AddJobs(map[string]*github.Job{
		"release": {
			Name:   String("create release"),
			RunsOn: String("ubuntu-latest"),
			Needs: &[]*string{
				String("golangci"),
				String("build"),
			},
			Steps: &[]*github.JobStep{
				{
					Name: String("Checkout"),
					Uses: String("actions/checkout@v2"),
				},
				{
					Name: String("Create Release"),
					Uses: String("go-semantic-release/action@v1"),
					Id:   String("create-release"),
					With: &map[string]interface{}{
						"github-token":             fmt.Sprintf("${{ secrets.%s }}", proj.GetGitHubToken()),
						"changelog-generator-opt":  "emojis=false",
						"force-bump-patch-version": true,
					},
				},
			},
		},
	})

	yaml, err := wf.ConvertToYAML()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/release.yml", pwd), yaml, 0o644)
	if err != nil {
		return err
	}

	wf2 := github.CreateWorkflow("Upload Binary")

	wf2.AddTrigger(github.Triggers{
		Release: &github.ReleaseOptions{
			Types: &[]*string{String("published")},
		},
	})

	wf2.AddJobs(map[string]*github.Job{
		"upload-binary": {
			Name:   String("upload binary"),
			RunsOn: String("ubuntu-latest"),
			Steps: &[]*github.JobStep{
				{
					Name: String("Checkout"),
					Uses: String("actions/checkout@v2"),
				},
				{
					Name: String("Upload binary"),
					Uses: String("wangyoucao577/go-release-action@v1.19"),
					With: &map[string]interface{}{
						"github_token": fmt.Sprintf("${{ secrets.%s }}", proj.GetGitHubToken()),
						"goos":         "linux",
						"goarch":       "amd64",
						"goversion":    proj.GetGoVersion(),
					},
				},
			},
		},
	})

	yaml2, err := wf2.ConvertToYAML()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/upload-binary.yml", pwd), yaml2, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) CreateBuildWorkflow() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.github/workflows", pwd), 0o755)
	if err != nil {
		return err
	}

	wf := github.CreateWorkflow("build")

	wf.AddTrigger(github.Triggers{
		PullRequest: &github.PullRequestOptions{},
	})

	wf, err = proj.setCommonJobs(wf)
	if err != nil {
		return err
	}

	err = wf.AddStep("build", &github.JobStep{
		Name: String("Commit and push changes (if changed)"),
		If:   String("steps.git_diff.outputs.has_changes"),
		Run: String(`git add . && git commit -m 'chore: self mutation && git push origin
HEAD:${{ github.event.pull_request.head.ref }}'`),
	})
	if err != nil {
		return err
	}

	err = wf.AddStep("build", &github.JobStep{
		Name: String("Update status check (if changed)"),
		If:   String("steps.git_diff.outputs.has_changes"),
		Run: String(`gh api -X POST /repos/${{ github.event.pull_request.head.repo.full_name
}}/check-runs -F name="build" -F head_sha="$(git rev-parse HEAD)" -F status="completed" -F conclusion="success`),
		Env: &map[string]*string{
			"GITHUB_TOKEN": String(fmt.Sprintf("${{ secrets.%s }}", proj.GetGitHubToken())),
		},
	})
	if err != nil {
		return err
	}

	err = wf.AddStep("build", &github.JobStep{
		Name: String("Cancel workflow (if changed)"),
		If:   String("steps.git_diff.outputs.has_changes"),
		Run: String(`gh api -X POST /repos/${{ github.event.pull_request.head.repo.full_name
}}/actions/runs/${{ github.run_id }}/cancel`),
		Env: &map[string]*string{
			"GITHUB_TOKEN": String(fmt.Sprintf("${{ secrets.%s }}", proj.GetGitHubToken())),
		},
	})
	if err != nil {
		return err
	}

	yaml, err := wf.ConvertToYAML()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/build.yml", pwd), yaml, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) GetName() string {
	if proj.Name == nil {
		return ""
	}
	return *proj.Name
}

func (proj *Project) GetDescription() string {
	if proj.Description == nil {
		return ""
	}
	return *proj.Description
}

func (proj *Project) GetRepository() string {
	if proj.Repository == nil {
		return ""
	}
	return *proj.Repository
}

func (proj *Project) GetAuthorName() string {
	if proj.AuthorName == nil {
		return ""
	}
	return *proj.AuthorName
}

func (proj *Project) GetAuthorEmail() string {
	if proj.AuthorEmail == nil {
		return ""
	}
	return *proj.AuthorEmail
}

func (proj *Project) GetGoVersion() string {
	if proj.GoVersion == nil {
		return "1.16"
	}
	return *proj.GoVersion
}

func (proj *Project) GetGitHubToken() string {
	if proj.GithubToken == nil {
		return "GITHUB_TOKEN"
	}
	return *proj.GithubToken
}

func (proj *Project) GetAuthorOrganization() string {
	if proj.AuthorOrganization == nil {
		return ""
	}
	return *proj.AuthorOrganization
}

func (proj *Project) IsRelease() bool {
	if proj.Release == nil {
		return false
	}
	return *proj.Release
}

func (proj *Project) GetDefaultReleaseBranch() string {
	if proj.DefaultReleaseBranch == nil {
		return "master"
	}
	return *proj.DefaultReleaseBranch
}

func (proj *Project) GetGitignore() []string {
	if proj.Gitignore == nil {
		return []string{}
	}
	return *proj.Gitignore
}

func (proj *Project) GetCodeOwners() []string {
	if proj.CodeOwners == nil {
		return []string{}
	}
	return *proj.CodeOwners
}

func (proj *Project) IsGoLinter() bool {
	if proj.GoLinter == nil {
		return false
	}
	return *proj.GoLinter
}

func (proj *Project) IsGoTest() bool {
	if proj.GoTest == nil {
		return true
	}
	return *proj.GoTest
}

func (proj *Project) GetGoTestArgs() []string {
	if proj.GoTestArgs == nil {
		return []string{}
	}
	return *proj.GoTestArgs
}

func (proj *Project) IsBuildWorkflow() bool {
	if proj.BuildWorkflow == nil {
		return false
	}
	return *proj.BuildWorkflow
}

func (proj *Project) IsIsGojen() bool {
	if proj.IsGojen == nil {
		return false
	}

	return *proj.IsGojen
}

func (proj *Project) GetGojenVersion() string {
	if proj.GojenVersion == nil {
		return "latest"
	}
	return *proj.GojenVersion
}

func (proj *Project) IsCreateReadme() bool {
	if proj.Readme == nil {
		return true
	}
	return *proj.Readme
}

func (proj *Project) IsCodeCov() bool {
	if proj.CodeCov == nil {
		return false
	}
	return *proj.CodeCov
}

func (proj *Project) GetTestEnvVars() []string {
	if proj.TestEnvVars == nil {
		return []string{}
	}
	return *proj.TestEnvVars
}

func (proj *Project) GetGoBuildArgs() []string {
	if proj.GoBuildArgs == nil {
		return []string{}
	}
	return *proj.GoBuildArgs
}

func (proj *Project) GetLicense() string {
	if proj.License == nil {
		return ""
	}
	return *proj.License
}

func (proj *Project) IsGoBuild() bool {
	if proj.GoBuild == nil {
		return true
	}
	return *proj.GoBuild
}

func (proj *Project) GetWorkflowEnv() *map[string]*string {
	if proj.WorkflowEnv == nil {
		return &map[string]*string{}
	}
	return proj.WorkflowEnv
}

func String(str string) *string {
	return &str
}

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

func StringSlice(s []string) *[]string {
	return &s
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
