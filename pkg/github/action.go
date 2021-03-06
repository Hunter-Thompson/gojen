package github

import (
	"errors"

	"gopkg.in/yaml.v2"
)

type IAction interface {
	// adds the job to the end of the action
	AddJobs(job map[string]*Job)
	// converts the action to yaml
	ConvertToYAML() ([]byte, error)
	// add triger to the action
	AddTrigger(trigger Triggers)
	// add steps to a job
	AddStep(jobName string, steps *JobStep) error
	// prepend steps to a job
	PrependStep(jobName string, steps *JobStep) error
	// append steps to a job
	AppendStep(jobName string, steps *JobStep) error
}

func CreateWorkflow(name string) IAction {
	return &Action{
		Name: name,
		On:   Triggers{},
		Jobs: make(map[string]*Job),
	}
}

// adds the job to the end of the action
func (a *Action) AddJobs(job map[string]*Job) {
	for k, v := range job {
		a.Jobs[k] = v
	}
}

func (a *Action) AddStep(jobName string, steps *JobStep) error {
	if a.Jobs[jobName] == nil {
		return errors.New("job does not exist")
	}
	*a.Jobs[jobName].Steps = append(*a.Jobs[jobName].Steps, steps)
	return nil
}

// add triger to the action
func (a *Action) AddTrigger(trigger Triggers) {
	a.On = trigger
}

// prepend steps to a job
func (a *Action) PrependStep(jobName string, steps *JobStep) error {
	if a.Jobs[jobName] == nil {
		return errors.New("job does not exist")
	}

	*a.Jobs[jobName].Steps = append([]*JobStep{steps}, *a.Jobs[jobName].Steps...)
	return nil
}

// append steps to a job
func (a *Action) AppendStep(jobName string, steps *JobStep) error {
	if a.Jobs[jobName] == nil {
		return errors.New("job does not exist")
	}

	*a.Jobs[jobName].Steps = append(*a.Jobs[jobName].Steps, steps)
	return nil
}

func (a *Action) ConvertToYAML() ([]byte, error) {
	b, err := yaml.Marshal(a)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// Check run options.
// Experimental.
type CheckRunOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

type Action struct {
	On   Triggers        `yaml:"on,omitempty,omitempty"`
	Name string          `yaml:"name,omitempty,omitempty"`
	Jobs map[string]*Job `yaml:"jobs,omitempty,omitempty"`
}

// Check suite options.
// Experimental.
type CheckSuiteOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Credentials to use to authenticate to Docker registries.
// Experimental.
type ContainerCredentials struct {
	// The password.
	// Experimental.
	Password *string `yaml:"password,omitempty"`
	// The username.
	// Experimental.
	Username *string `yaml:"username,omitempty"`
}

// Options petaining to container environments.
// Experimental.
type ContainerOptions struct {
	// The Docker image to use as the container to run the action.
	//
	// The value can
	// be the Docker Hub image name or a registry name.
	// Experimental.
	Image *string `yaml:"image,omitempty"`
	// f the image's container registry requires authentication to pull the image, you can use credentials to set a map of the username and password.
	//
	// The credentials are the same values that you would provide to the docker
	// login command.
	// Experimental.
	Credentials *ContainerCredentials `yaml:"credentials,omitempty"`
	// Sets a map of environment variables in the container.
	// Experimental.
	Env *map[string]*string `yaml:"env,omitempty"`
	// Additional Docker container resource options.
	// See: https://docs.docker.com/engine/reference/commandline/create/#options
	//
	// Experimental.
	Options *[]*string `yaml:"options,omitempty"`
	// Sets an array of ports to expose on the container.
	// Experimental.
	Ports *[]*float64 `yaml:"ports,omitempty"`
	// Sets an array of volumes for the container to use.
	//
	// You can use volumes to
	// share data between services or other steps in a job. You can specify
	// named Docker volumes, anonymous Docker volumes, or bind mounts on the
	// host.
	//
	// To specify a volume, you specify the source and destination path:
	// `<source>:<destinationPath>`.
	// Experimental.
	Volumes *[]*string `yaml:"volumes,omitempty"`
}

// The Create event accepts no options.
// Experimental.
type CreateOptions struct {
}

// CRON schedule options.
// Experimental.
type CronScheduleOptions struct {
	// See: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html#tag_20_25_07
	//
	// Experimental.
	Cron *string `yaml:"cron,omitempty"`
}

// The Delete event accepts no options.
// Experimental.
type DeleteOptions struct {
}

// The Deployment event accepts no options.
// Experimental.
type DeploymentOptions struct {
}

// The Deployment status event accepts no options.
// Experimental.
type DeploymentStatusOptions struct {
}

// The Fork event accepts no options.
// Experimental.
type ForkOptions struct {
}

// The Gollum event accepts no options.
// Experimental.
type GollumOptions struct {
}

// Issue comment options.
// Experimental.
type IssueCommentOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Issues options.
// Experimental.
type IssuesOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// A GitHub Workflow job definition.
// Experimental.
type Job struct {
	// You can modify the default permissions granted to the GITHUB_TOKEN, adding or removing access as required, so that you only allow the minimum required access.
	//
	// Use `{ contents: READ }` if your job only needs to clone code.
	//
	// This is intentionally a required field since it is required in order to
	// allow workflows to run in GitHub repositories with restricted default
	// access.
	// See: https://docs.github.com/en/actions/reference/authentication-in-a-workflow#permissions-for-the-github_token
	//
	// Experimental.
	Permissions *JobPermissions `yaml:"permissions,omitempty"`
	// The type of machine to run the job on.
	//
	// The machine can be either a
	// GitHub-hosted runner or a self-hosted runner.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	RunsOn *string `yaml:"runs-on,omitempty"`
	// A job contains a sequence of tasks called steps.
	//
	// Steps can run commands,
	// run setup tasks, or run an action in your repository, a public repository,
	// or an action published in a Docker registry. Not all steps run actions,
	// but all actions run as a step. Each step runs in its own process in the
	// runner environment and has access to the workspace and filesystem.
	// Because steps run in their own process, changes to environment variables
	// are not preserved between steps. GitHub provides built-in steps to set up
	// and complete a job.
	// Experimental.
	Steps *[]*JobStep `yaml:"steps,omitempty"`
	// Concurrency ensures that only a single job or workflow using the same concurrency group will run at a time.
	//
	// A concurrency group can be any
	// string or expression. The expression can use any context except for the
	// secrets context.
	// Experimental.
	Concurrency interface{} `yaml:"concurrency,omitempty"`
	// A container to run any steps in a job that don't already specify a container.
	//
	// If you have steps that use both script and container actions,
	// the container actions will run as sibling containers on the same network
	// with the same volume mounts.
	// Experimental.
	Container *ContainerOptions `yaml:"container,omitempty"`
	// Prevents a workflow run from failing when a job fails.
	//
	// Set to true to
	// allow a workflow run to pass when this job fails.
	// Experimental.
	ContinueOnError *bool `yaml:"continueOnError,omitempty"`
	// A map of default settings that will apply to all steps in the job.
	//
	// You
	// can also set default settings for the entire workflow.
	// Experimental.
	Defaults *JobDefaults `yaml:"defaults,omitempty"`
	// A map of environment variables that are available to all steps in the job.
	//
	// You can also set environment variables for the entire workflow or an
	// individual step.
	// Experimental.
	Env *map[string]*string `yaml:"env,omitempty"`
	// The environment that the job references.
	//
	// All environment protection rules
	// must pass before a job referencing the environment is sent to a runner.
	// See: https://docs.github.com/en/actions/reference/environments
	//
	// Experimental.
	Environment interface{} `yaml:"environment,omitempty"`
	// You can use the if conditional to prevent a job from running unless a condition is met.
	//
	// You can use any supported context and expression to
	// create a conditional.
	// Experimental.
	If *string `yaml:"if,omitempty"`
	// The name of the job displayed on GitHub.
	// Experimental.
	Name *string `yaml:"name,omitempty"`
	// Identifies any jobs that must complete successfully before this job will run.
	//
	// It can be a string or array of strings. If a job fails, all jobs
	// that need it are skipped unless the jobs use a conditional expression
	// that causes the job to continue.
	// Experimental.
	Needs *[]*string `yaml:"needs,omitempty"`
	// A map of outputs for a job.
	//
	// Job outputs are available to all downstream
	// jobs that depend on this job.
	// Experimental.
	Outputs *map[string]*JobStepOutput `yaml:"outputs,omitempty"`
	// Used to host service containers for a job in a workflow.
	//
	// Service
	// containers are useful for creating databases or cache services like Redis.
	// The runner automatically creates a Docker network and manages the life
	// cycle of the service containers.
	// Experimental.
	Services *map[string]*ContainerOptions `yaml:"services,omitempty"`
	// A strategy creates a build matrix for your jobs.
	//
	// You can define different
	// variations to run each job in.
	// Experimental.
	Strategy *JobStrategy `yaml:"strategy,omitempty"`
	// The maximum number of minutes to let a job run before GitHub automatically cancels it.
	// Experimental.
	TimeoutMinutes *float64 `yaml:"timeout-minutes,omitempty"`
}

// Default settings for all steps in the job.
// Experimental.
type JobDefaults struct {
	// Default run settings.
	// Experimental.
	Run *RunSettings `yaml:"run,omitempty"`
}

// A job matrix.
// Experimental.
type JobMatrix struct {
	// Each option you define in the matrix has a key and value.
	//
	// The keys you
	// define become properties in the matrix context and you can reference the
	// property in other areas of your workflow file. For example, if you define
	// the key os that contains an array of operating systems, you can use the
	// matrix.os property as the value of the runs-on keyword to create a job
	// for each operating system.
	// Experimental.
	Domain *map[string]*[]*string `yaml:"domain,omitempty"`
	// You can remove a specific configurations defined in the build matrix using the exclude option.
	//
	// Using exclude removes a job defined by the
	// build matrix.
	// Experimental.
	Exclude *[]*map[string]*string `yaml:"exclude,omitempty"`
	// You can add additional configuration options to a build matrix job that already exists.
	//
	// For example, if you want to use a specific version of npm
	// when the job that uses windows-latest and version 8 of node runs, you can
	// use include to specify that additional option.
	// Experimental.
	Include *[]*map[string]*string `yaml:"include,omitempty"`
}

// Access level for workflow permission scopes.
// Experimental.
type JobPermission string

const (
	JobPermission_READ  JobPermission = "READ"
	JobPermission_WRITE JobPermission = "WRITE"
	JobPermission_NONE  JobPermission = "NONE"
)

// The available scopes and access values for workflow permissions.
//
// If you
// specify the access for any of these scopes, all those that are not
// specified are set to `JobPermission.NONE`, instead of the default behavior
// when none is specified.
// Experimental.
type JobPermissions struct {
	// Experimental.
	Actions JobPermission `yaml:"actions,omitempty"`
	// Experimental.
	Checks JobPermission `yaml:"checks,omitempty"`
	// Experimental.
	Contents JobPermission `yaml:"contents,omitempty"`
	// Experimental.
	Deployments JobPermission `yaml:"deployments,omitempty"`
	// Experimental.
	Issues JobPermission `yaml:"issues,omitempty"`
	// Experimental.
	Packages JobPermission `yaml:"packages,omitempty"`
	// Experimental.
	PullRequests JobPermission `yaml:"pull_requests,omitempty"`
	// Experimental.
	RepositoryProjects JobPermission `yaml:"repository_projects,omitempty"`
	// Experimental.
	SecurityEvents JobPermission `yaml:"security_events,omitempty"`
	// Experimental.
	Statuses JobPermission `yaml:"statuses,omitempty"`
}

// A job step.
// Experimental.
type JobStep struct {
	// Prevents a job from failing when a step fails.
	//
	// Set to true to allow a job
	// to pass when this step fails.
	// Experimental.
	ContinueOnError *bool `yaml:"continue-on-error,omitempty"`
	// Sets environment variables for steps to use in the runner environment.
	//
	// You can also set environment variables for the entire workflow or a job.
	// Experimental.
	Env *map[string]*string `yaml:"env,omitempty"`
	// A unique identifier for the step.
	//
	// You can use the id to reference the
	// step in contexts.
	// Experimental.
	Id *string `yaml:"id,omitempty"`
	// You can use the if conditional to prevent a job from running unless a condition is met.
	//
	// You can use any supported context and expression to
	// create a conditional.
	// Experimental.
	If *string `yaml:"if,omitempty"`
	// A name for your step to display on GitHub.
	// Experimental.
	Name *string `yaml:"name,omitempty"`
	// Runs command-line programs using the operating system's shell.
	//
	// If you do
	// not provide a name, the step name will default to the text specified in
	// the run command.
	// Experimental.
	Run *string `yaml:"run,omitempty"`
	// The maximum number of minutes to run the step before killing the process.
	// Experimental.
	TimeoutMinutes *float64 `yaml:"timeout-minutes,omitempty"`
	// Selects an action to run as part of a step in your job.
	//
	// An action is a
	// reusable unit of code. You can use an action defined in the same
	// repository as the workflow, a public repository, or in a published Docker
	// container image.
	// Experimental.
	Uses *string `yaml:"uses,omitempty"`
	// A map of the input parameters defined by the action.
	//
	// Each input parameter
	// is a key/value pair. Input parameters are set as environment variables.
	// The variable is prefixed with INPUT_ and converted to upper case.
	// Experimental.
	With *map[string]interface{} `yaml:"with,omitempty"`
}

// An output binding for a job.
// Experimental.
type JobStepOutput struct {
	// The name of the job output that is being bound.
	// Experimental.
	OutputName *string `yaml:"output-name,omitempty"`
	// The ID of the step that exposes the output.
	// Experimental.
	StepId *string `yaml:"step-id,omitempty"`
}

// A strategy creates a build matrix for your jobs.
//
// You can define different
// variations to run each job in.
// Experimental.
type JobStrategy struct {
	// When set to true, GitHub cancels all in-progress jobs if any matrix job fails.
	//
	// Default: true
	// Experimental.
	FailFast *bool `yaml:"fail-fast,omitempty"`
	// You can define a matrix of different job configurations.
	//
	// A matrix allows
	// you to create multiple jobs by performing variable substitution in a
	// single job definition. For example, you can use a matrix to create jobs
	// for more than one supported version of a programming language, operating
	// system, or tool. A matrix reuses the job's configuration and creates a
	// job for each matrix you configure.
	//
	// A job matrix can generate a maximum of 256 jobs per workflow run. This
	// limit also applies to self-hosted runners.
	// Experimental.
	Matrix *JobMatrix `yaml:"matrix,omitempty"`
	// The maximum number of jobs that can run simultaneously when using a matrix job strategy.
	//
	// By default, GitHub will maximize the number of jobs
	// run in parallel depending on the available runners on GitHub-hosted
	// virtual machines.
	// Experimental.
	MaxParallel *float64 `yaml:"max-parallel,omitempty"`
}

// label options.
// Experimental.
type LabelOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Milestone options.
// Experimental.
type MilestoneOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// The Page build event accepts no options.
// Experimental.
type PageBuildOptions struct {
}

// Project card options.
// Experimental.
type ProjectCardOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Probject column options.
// Experimental.
type ProjectColumnOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Project options.
// Experimental.
type ProjectOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// The Public event accepts no options.
// Experimental.
type PublicOptions struct {
}

// Pull request options.
// Experimental.
type PullRequestOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Pull request review comment options.
// Experimental.
type PullRequestReviewCommentOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Pull request review options.
// Experimental.
type PullRequestReviewOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Pull request target options.
// Experimental.
type PullRequestTargetOptions struct {
	// When using the push and pull_request events, you can configure a workflow to run on specific branches or tags.
	//
	// For a pull_request event, only
	// branches and tags on the base are evaluated. If you define only tags or
	// only branches, the workflow won't run for events affecting the undefined
	// Git ref.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Branches *[]*string `yaml:"branches,omitempty"`
	// When using the push and pull_request events, you can configure a workflow to run when at least one file does not match paths-ignore or at least one modified file matches the configured paths.
	//
	// Path filters are not
	// evaluated for pushes to tags.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Paths *[]*string `yaml:"paths,omitempty"`
	// When using the push and pull_request events, you can configure a workflow to run on specific branches or tags.
	//
	// For a pull_request event, only
	// branches and tags on the base are evaluated. If you define only tags or
	// only branches, the workflow won't run for events affecting the undefined
	// Git ref.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Tags *[]*string `yaml:"tags,omitempty"`
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Options for push-like events.
// Experimental.
type PushOptions struct {
	// When using the push and pull_request events, you can configure a workflow to run on specific branches or tags.
	//
	// For a pull_request event, only
	// branches and tags on the base are evaluated. If you define only tags or
	// only branches, the workflow won't run for events affecting the undefined
	// Git ref.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Branches *[]*string `yaml:"branches,omitempty"`
	// When using the push and pull_request events, you can configure a workflow to run when at least one file does not match paths-ignore or at least one modified file matches the configured paths.
	//
	// Path filters are not
	// evaluated for pushes to tags.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Paths *[]*string `yaml:"paths,omitempty"`
	// When using the push and pull_request events, you can configure a workflow to run on specific branches or tags.
	//
	// For a pull_request event, only
	// branches and tags on the base are evaluated. If you define only tags or
	// only branches, the workflow won't run for events affecting the undefined
	// Git ref.
	// See: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	//
	// Experimental.
	Tags *[]*string `yaml:"tags,omitempty"`
}

// Registry package options.
// Experimental.
type RegistryPackageOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Release options.
// Experimental.
type ReleaseOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Repository dispatch options.
// Experimental.
type RepositoryDispatchOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// Run settings for a job.
// Experimental.
type RunSettings struct {
	// Which shell to use for running the step.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Shell *string `yaml:"shell,omitempty"`
	// Working directory to use when running the step.
	// Experimental.
	WorkingDirectory *string `yaml:"working-directory,omitempty"`
}

// The Status event accepts no options.
// Experimental.
type StatusOptions struct {
}

// The set of available triggers for GitHub Workflows.
// See: https://docs.github.com/en/actions/reference/events-that-trigger-workflows
//
// Experimental.
type Triggers struct {
	// Runs your workflow anytime the check_run event occurs.
	// Experimental.
	CheckRun *CheckRunOptions `yaml:"check_run,omitempty"`
	// Runs your workflow anytime the check_suite event occurs.
	// Experimental.
	CheckSuite *CheckSuiteOptions `yaml:"check-suite,omitempty"`
	// Runs your workflow anytime someone creates a branch or tag, which triggers the create event.
	// Experimental.
	Create *CreateOptions `yaml:"create,omitempty"`
	// Runs your workflow anytime someone deletes a branch or tag, which triggers the delete event.
	// Experimental.
	Delete *DeleteOptions `yaml:"delete,omitempty"`
	// Runs your workflow anytime someone creates a deployment, which triggers the deployment event.
	//
	// Deployments created with a commit SHA may not have
	// a Git ref.
	// Experimental.
	Deployment *DeploymentOptions `yaml:"deployment,omitempty"`
	// Runs your workflow anytime a third party provides a deployment status, which triggers the deployment_status event.
	//
	// Deployments created with a
	// commit SHA may not have a Git ref.
	// Experimental.
	DeploymentStatus *DeploymentStatusOptions `yaml:"deployment_status,omitempty"`
	// Runs your workflow anytime when someone forks a repository, which triggers the fork event.
	// Experimental.
	Fork *ForkOptions `yaml:"fork,omitempty"`
	// Runs your workflow when someone creates or updates a Wiki page, which triggers the gollum event.
	// Experimental.
	Gollum *GollumOptions `yaml:"gollum,omitempty"`
	// Runs your workflow anytime the issue_comment event occurs.
	// Experimental.
	IssueComment *IssueCommentOptions `yaml:"issue_comment,omitempty"`
	// Runs your workflow anytime the issues event occurs.
	// Experimental.
	Issues *IssuesOptions `yaml:"issues,omitempty"`
	// Runs your workflow anytime the label event occurs.
	// Experimental.
	Label *LabelOptions `yaml:"label,omitempty"`
	// Runs your workflow anytime the milestone event occurs.
	// Experimental.
	Milestone *MilestoneOptions `yaml:"milestone,omitempty"`
	// Runs your workflow anytime someone pushes to a GitHub Pages-enabled branch, which triggers the page_build event.
	// Experimental.
	PageBuild *PageBuildOptions `yaml:"page_build,omitempty"`
	// Runs your workflow anytime the project event occurs.
	// Experimental.
	Project *ProjectOptions `yaml:"project,omitempty"`
	// Runs your workflow anytime the project_card event occurs.
	// Experimental.
	ProjectCard *ProjectCardOptions `yaml:"project_card,omitempty"`
	// Runs your workflow anytime the project_column event occurs.
	// Experimental.
	ProjectColumn *ProjectColumnOptions `yaml:"project_column,omitempty"`
	// Runs your workflow anytime someone makes a private repository public, which triggers the public event.
	// Experimental.
	Public *PublicOptions `yaml:"public,omitempty"`
	// Runs your workflow anytime the pull_request event occurs.
	// Experimental.
	PullRequest *PullRequestOptions `yaml:"pull_request,omitempty"`
	// Runs your workflow anytime the pull_request_review event occurs.
	// Experimental.
	PullRequestReview *PullRequestReviewOptions `yaml:"pull_request_review,omitempty"`
	// Runs your workflow anytime a comment on a pull request's unified diff is modified, which triggers the pull_request_review_comment event.
	// Experimental.
	PullRequestReviewComment *PullRequestReviewCommentOptions `yaml:"pull_request_review_comment,omitempty"`
	// This event runs in the context of the base of the pull request, rather than in the merge commit as the pull_request event does.
	//
	// This prevents
	// executing unsafe workflow code from the head of the pull request that
	// could alter your repository or steal any secrets you use in your workflow.
	// This event allows you to do things like create workflows that label and
	// comment on pull requests based on the contents of the event payload.
	//
	// WARNING: The `pull_request_target` event is granted read/write repository
	// token and can access secrets, even when it is triggered from a fork.
	// Although the workflow runs in the context of the base of the pull request,
	// you should make sure that you do not check out, build, or run untrusted
	// code from the pull request with this event. Additionally, any caches
	// share the same scope as the base branch, and to help prevent cache
	// poisoning, you should not save the cache if there is a possibility that
	// the cache contents were altered.
	// See: https://securitylab.github.com/research/github-actions-preventing-pwn-requests
	//
	// Experimental.
	PullRequestTarget *PullRequestTargetOptions `yaml:"pull-request-target,omitempty"`
	// Runs your workflow when someone pushes to a repository branch, which triggers the push event.
	// Experimental.
	Push *PushOptions `yaml:"push,omitempty"`
	// Runs your workflow anytime a package is published or updated.
	// Experimental.
	RegistryPackage *RegistryPackageOptions `yaml:"registry-package,omitempty"`
	// Runs your workflow anytime the release event occurs.
	// Experimental.
	Release *ReleaseOptions `yaml:"release,omitempty"`
	// You can use the GitHub API to trigger a webhook event called repository_dispatch when you want to trigger a workflow for activity that happens outside of GitHub.
	// Experimental.
	RepositoryDispatch *RepositoryDispatchOptions `yaml:"repository-dispatch,omitempty"`
	// You can schedule a workflow to run at specific UTC times using POSIX cron syntax.
	//
	// Scheduled workflows run on the latest commit on the default or
	// base branch. The shortest interval you can run scheduled workflows is
	// once every 5 minutes.
	// See: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html#tag_20_25_07
	//
	// Experimental.
	Schedule *[]*CronScheduleOptions `yaml:"schedule,omitempty"`
	// Runs your workflow anytime the status of a Git commit changes, which triggers the status event.
	// Experimental.
	Status *StatusOptions `yaml:"status,omitempty"`
	// Runs your workflow anytime the watch event occurs.
	// Experimental.
	Watch *WatchOptions `yaml:"watch,omitempty"`
	// You can configure custom-defined input properties, default input values, and required inputs for the event directly in your workflow.
	//
	// When the
	// workflow runs, you can access the input values in the github.event.inputs
	// context.
	// Experimental.
	WorkflowDispatch *WorkflowDispatchOptions `yaml:"workflow-dispatch,omitempty"`
	// This event occurs when a workflow run is requested or completed, and allows you to execute a workflow based on the finished result of another workflow.
	//
	// A workflow run is triggered regardless of the result of the
	// previous workflow.
	// Experimental.
	WorkflowRun *WorkflowRunOptions `yaml:"workflow-run,omitempty"`
}

// Watch options.
// Experimental.
type WatchOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}

// The Workflow dispatch event accepts no options.
// Experimental.
type WorkflowDispatchOptions struct {
}

// Workflow run options.
// Experimental.
type WorkflowRunOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `yaml:"types,omitempty"`
}
