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
)

type IProject interface {
	WriteConfig() error
	SetupProject() error
	SetGitignore() error
	RunTest() error
	RunLinter() error

	GetName() string
	GetDescription() string
	GetRepository() string
	GetAuthorName() string
	GetAuthorEmail() string
	GetAuthorOrganization() string
	IsLicensed() bool
	IsRelease() bool
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
}

type Project struct {
	Name        *string `yaml:"name" json:"name"`
	Description *string `yaml:"description" json:"description"`
	Repository  *string `yaml:"repository"  json:"repository"`
	GoVersion   *string `yaml:"goVersion" json:"goVersion"`

	AuthorName         *string `yaml:"authorName" json:"authorName"`
	AuthorEmail        *string `yaml:"authorEmail" json:"authorEmail"`
	AuthorOrganization *string `yaml:"authorOrganization" json:"authorOrganization"`

	Licensed     *bool   `yaml:"licensed" json:"licensed"`
	GojenVersion *string `yaml:"gojenVersion" json:"gojenVersion"`

	Release              *bool   `yaml:"release" json:"release"`
	BuildWorkflow        *bool   `yaml:"buildWorkflow" json:"buildWorkflow"`
	GithubToken          *string `yaml:"githubToken" json:"githubToken"`
	DefaultReleaseBranch *string `yaml:"defaultReleaseBranch" json:"defaultReleaseBranch"`
	IsGojen              *bool   `yaml:"isGojen" json:"isGojen"`

	Gitignore  *[]string `yaml:"gitignore" json:"gitignore"`
	CodeOwners *[]string `yaml:"codeOwners" json:"codeOwners"`

	GoLinter   *bool     `yaml:"goLinter" json:"goLinter"`
	GoTest     *bool     `yaml:"goTest" json:"goTest"`
	GoTestArgs *[]string `yaml:"goTestArgs" json:"goTestArgs"`
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

	err = ioutil.WriteFile(cfgPath, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) SetupProject() error {

	err := proj.SetGitignore()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(proj.CodeOwners, &[]string{}) {
		err = proj.SetCodeOwners()
		if err != nil {
			return err
		}
	}

	if proj.IsRelease() {
		err = proj.CreateReleaseWorkflow()
		if err != nil {
			return err
		}
	}

	if proj.IsBuildWorkflow() {
		err = proj.CreateBuildWorkflow()
		if err != nil {
			return err
		}
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
		err := ioutil.WriteFile(pwd+"/main.go", []byte(c+"\n"), 0644)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if _, err := os.Stat(pwd + "/go.mod"); errors.Is(err, os.ErrNotExist) {
		out, err := modInit.CombinedOutput()
		if err != nil {
			fmt.Println("running go mod init failed")
			return errors.New(string(out))
		}
		fmt.Print(string(out))
	} else {
		fmt.Println("go.mod already exists, moving on ...")
	}

	if err != nil {
		return err
	}

	fmt.Println("running go mod vendor")
	out, err := vendor.CombinedOutput()
	if err != nil {
		fmt.Println("running go mod vendor failed")
		return errors.New(string(out))
	}

	fmt.Println("running go mod tidy")
	out, err = tidy.CombinedOutput()
	if err != nil {
		fmt.Println("running go mod tidy failed")
		return errors.New(string(out))
	}

	fmt.Println("running go fmt")
	out, err = gofmt.CombinedOutput()
	if err != nil {
		fmt.Println("running go fmt failed")
		return errors.New(string(out))
	}

	if proj.IsGoLinter() {
		err = proj.RunLinter()
		if err != nil {
			return err
		}
	}

	if proj.IsGoTest() {
		err = proj.RunTest()
		if err != nil {
			return err
		}
	}

	return nil
}

func (proj *Project) RunTest() error {

	fmt.Println("running go test")

	out, err := exec.Command("go", append([]string{"test", "-v"}, *proj.GoTestArgs...)...).CombinedOutput()
	if err != nil {
		fmt.Println("running go test failed")
		return errors.New(string(out))
	}

	fmt.Print(string(out))
	fmt.Println("go test passed")
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

	err = ioutil.WriteFile(gitignorePath, []byte(contents), 0644)
	if err != nil {
		return err
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
		err = os.MkdirAll(fmt.Sprintf("%s/.github", pwd), 0755)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("codeowners file exists, moving on ...")
	}

	contents := strings.Join(*proj.CodeOwners, "\n")

	err = ioutil.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) RunLinter() error {

	fmt.Println("running go linter")

	out, err := exec.Command("golangci-lint", "run").CombinedOutput()
	if err != nil {
		fmt.Println("running golint failed")
		return errors.New(string(out))
	}

	fmt.Print(string(out))
	fmt.Println("go linter passed")
	return nil

}

func (proj *Project) CreateReleaseWorkflow() error {

	fmt.Println("creating release workflow")

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.github/workflows", pwd), 0755)
	if err != nil {
		return err
	}

	var gojenCommand string

	if proj.IsIsGojen() {
		gojenCommand = `- name: build and run gojen
	    run: "go build && gojen"
		`
	} else {
		gojenCommand = `- name: Install gojen
      run: go install github.com/Hunter-Thompson/gojen
	  - name: Run gojen
	    run: gojen 
		`
	}

	c := fmt.Sprintf(`name: Release
on:
  push:
	branches:
	  - %s

jobs:
  build:
    runs-on: ubuntu:latest
    name: Release
    steps:
    - uses: actions/checkout@v2
    - name: Setup go 
      uses: actions/setup-go@v2
      with:
        go-version: %s
		%s
  release:
    needs:
    - build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: go-semantic-release/action@v1
        id: semrel
        with:
          github-token: ${{ secrets.%s }}
          changelog-generator-opt: "emojis=false"
          force-bump-patch-version: true`, proj.GetDefaultReleaseBranch(), proj.GetGoVersion(), gojenCommand, proj.GetGitHubToken())

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/release.yml", pwd), []byte(c), 0644)
	if err != nil {
		return err
	}

	c = fmt.Sprintf(`name: upload binary

on:
  release:
    types: [published]

jobs:
  release:
    name: release 
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: wangyoucao577/go-release-action@v1.19
      with:
        github_token: ${{ secrets.%s }}
        goos: linux
        goarch: amd64
        goversion: %s`, proj.GetGitHubToken(), proj.GetGoVersion())

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/upload-binary.yml", pwd), []byte(c), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (proj *Project) CreateBuildWorkflow() error {

	fmt.Println("creating build workflow")

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.github/workflows", pwd), 0755)
	if err != nil {
		return err
	}

	var gojenCommand string

	if proj.IsIsGojen() {
		gojenCommand = `- name: build and run gojen
	    run: "go build && gojen"`
	} else {
		gojenCommand = fmt.Sprintf(`- name: Install gojen
      run: go install github.com/Hunter-Thompson/gojen@%s
	  - name: Run gojen
	    run: gojen `, proj.GetGojenVersion())
	}

	c := fmt.Sprintf(`name: Build
on:
  pull_request: {}

jobs:
  build:
	  runs-on: ubuntu:latest
	  name: Build
	  steps:
	  - uses: actions/checkout@v2
        with:
          ref: ${{ github.event.pull_request.head.ref }}
          repository: ${{ github.event.pull_request.head.repo.full_name }}
	  - name: Setup go 
	    uses: actions/setup-go@v2
	    with:
	  	go-version: %s
	  %s
    - name: Check for changes
      id: git_diff
      run: git diff --exit-code || echo "::set-output name=has_changes::true"
    - if: steps.git_diff.outputs.has_changes
      name: Commit and push changes (if changed)
      run: 'git add . && git commit -m "chore: self mutation" && git push origin
        HEAD:${{ github.event.pull_request.head.ref }}'
    - if: steps.git_diff.outputs.has_changes
      name: Update status check (if changed)
      run: gh api -X POST /repos/${{ github.event.pull_request.head.repo.full_name
        }}/check-runs -F name="build" -F head_sha="$(git rev-parse HEAD)" -F
        status="completed" -F conclusion="success"
      env:
        GITHUB_TOKEN: ${{ secrets.%s }}
    - if: steps.git_diff.outputs.has_changes
      name: Cancel workflow (if changed)
      run: gh api -X POST /repos/${{ github.event.pull_request.head.repo.full_name
        }}/actions/runs/${{ github.run_id }}/cancel
      env:
        GITHUB_TOKEN: ${{ secrets.%s }}`, proj.GetGoVersion(), gojenCommand, proj.GetGitHubToken(), proj.GetGitHubToken())

	err = ioutil.WriteFile(fmt.Sprintf("%s/.github/workflows/build.yml", pwd), []byte(c), 0644)
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

func (proj *Project) IsLicensed() bool {
	if proj.Licensed == nil {
		return false
	}
	return *proj.Licensed
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
		return "0.1.0"
	}
	return *proj.GojenVersion
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
