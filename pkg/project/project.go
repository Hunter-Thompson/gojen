package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
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
}

type Project struct {
	Name        *string `yaml:"name" json:"name"`
	Description *string `yaml:"description" json:"description"`
	Repository  *string `yaml:"repository"  json:"repository"`

	AuthorName         *string `yaml:"authorName" json:"authorName"`
	AuthorEmail        *string `yaml:"authorEmail" json:"authorEmail"`
	AuthorOrganization *string `yaml:"authorOrganization" json:"authorOrganization"`

	Licensed *bool `yaml:"licensed" json:"licensed"`

	Release              *bool   `yaml:"release" json:"release"`
	DefaultReleaseBranch *string `yaml:"defaultReleaseBranch" json:"defaultReleaseBranch"`

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
	if *proj.Name == "" {
		return errors.New("name is missing in config")
	}

	if *proj.Repository == "" {
		return errors.New("repository is missing in config")
	}

	if *proj.AuthorName == "" {
		return errors.New("authorname is missing in config")
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

	err = yaml.Unmarshal(b, &proj)
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

	// err := exec.Command("go", "get", "-u", "-v", "github.com/golang/lint/golint").Run()

	err := proj.SetGitignore()
	if err != nil {
		return err
	}

	err = proj.SetCodeOwners()
	if err != nil {
		return err
	}

	modInit := exec.Command("go", "mod", "init", proj.GetRepository())
	vendor := exec.Command("go", "mod", "vendor")
	tidy := exec.Command("go", "mod", "tidy")

	pwd, err := os.Getwd()
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
	err = vendor.Run()
	if err != nil {
		return err
	}

	fmt.Println("running go mod tidy")
	err = tidy.Run()
	if err != nil {
		return err
	}

	err = proj.RunLinter()
	if err != nil {
		return err
	}

	err = proj.RunTest()
	if err != nil {
		return err
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

func (proj *Project) GetName() string {
	return *proj.Name
}

func (proj *Project) GetDescription() string {
	return *proj.Description
}

func (proj *Project) GetRepository() string {
	return *proj.Repository
}

func (proj *Project) GetAuthorName() string {
	return *proj.AuthorName
}

func (proj *Project) GetAuthorEmail() string {
	return *proj.AuthorEmail
}

func (proj *Project) GetAuthorOrganization() string {
	return *proj.AuthorOrganization
}

func (proj *Project) IsLicensed() bool {
	return *proj.Licensed
}

func (proj *Project) IsRelease() bool {
	return *proj.Release
}

func (proj *Project) GetDefaultReleaseBranch() string {
	return *proj.DefaultReleaseBranch
}

func (proj *Project) GetGitignore() []string {
	return *proj.Gitignore
}

func (proj *Project) GetCodeOwners() []string {
	return *proj.CodeOwners
}

func (proj *Project) IsGoLinter() bool {
	return *proj.GoLinter
}

func (proj *Project) IsGoTest() bool {
	return *proj.GoTest
}

func (proj *Project) GetGoTestArgs() []string {
	return *proj.GoTestArgs
}
