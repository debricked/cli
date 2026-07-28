package scan

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/debricked/cli/internal/callgraph"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/ci"
	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/debug"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/fingerprint"
	"github.com/debricked/cli/internal/git"
	"github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/report/sbom"
	"github.com/debricked/cli/internal/resolution"
	"github.com/debricked/cli/internal/tui"
	"github.com/debricked/cli/internal/upload"
	"github.com/fatih/color"
)

var (
	BadOptsErr      = errors.New("failed to type case IOptions")
	FailPipelineErr = errors.New("")
	LongQueueErr    = errors.New("progress polling terminated due to long scan times")
)

type IScanner interface {
	Scan(o IOptions) error
}

type IOptions interface{}

type DebrickedScanner struct {
	client      *client.IDebClient
	finder      file.IFinder
	uploader    *upload.IUploader
	ciService   ci.IService
	resolver    resolution.IResolver
	fingerprint fingerprint.IFingerprint
	callgraph   callgraph.IGenerator
}

type DebrickedOptions struct {
	Path                        string
	Resolve                     bool
	Fingerprint                 bool
	CallGraph                   bool
	JavaCallgraphEngine         string
	SBOM                        string
	SBOMOutput                  string
	Exclusions                  []string
	Inclusions                  []string
	Verbose                     bool
	Debug                       bool
	Regenerate                  int
	VersionHint                 bool
	RepositoryName              string
	CommitName                  string
	GenerateCommitName          bool
	BranchName                  string
	CommitAuthor                string
	RepositoryUrl               string
	IntegrationName             string
	JsonFilePath                string
	NpmPreferred                bool
	PassOnTimeOut               bool
	CallGraphUploadTimeout      int
	CallGraphGenerateTimeout    int
	MinFingerprintContentLength int
	TagCommitAsRelease          bool
	Experimental                bool
	Version                     string
	ResolutionStrictness        resolution.StrictnessLevel
}

func NewDebrickedScanner(
	c *client.IDebClient,
	finder file.IFinder,
	uploader upload.IUploader,
	ciService ci.IService,
	resolver resolution.IResolver,
	fingerprint fingerprint.IFingerprint,
	callgraph callgraph.IGenerator,
) *DebrickedScanner {
	return &DebrickedScanner{
		c,
		finder,
		&uploader,
		ciService,
		resolver,
		fingerprint,
		callgraph,
	}
}

func (dScanner *DebrickedScanner) Scan(o IOptions) error {
	dOptions, ok := o.(DebrickedOptions)
	if !ok {
		return BadOptsErr
	}
	debug.Log("Options initialized, finding CI service...", dOptions.Debug)

	e, _ := dScanner.ciService.Find()

	debug.Log("Mapping environment variables...", dOptions.Debug)
	MapEnvToOptions(&dOptions, e)
	UpdatedEmptyCommitName(&dOptions)

	if err := SetWorkingDirectory(&dOptions); err != nil {
		return err
	}

	debug.Log("Setting up git objects...", dOptions.Debug)
	gitMetaObject, err := git.NewMetaObject(
		dOptions.Path,
		dOptions.RepositoryName,
		dOptions.CommitName,
		dOptions.BranchName,
		dOptions.CommitAuthor,
		dOptions.RepositoryUrl,
	)
	if err != nil {
		return err
	}

	debug.Log("Running scanResolve...", dOptions.Debug)
	resolutionErr := dScanner.scanResolve(dOptions)
	if isFatalResolutionErr(resolutionErr) {
		return resolutionErr
	}

	debug.Log("Running scan with initialized scanner...", dOptions.Debug)
	result, err := dScanner.scan(dOptions, *gitMetaObject)
	if err != nil {
		return dScanner.handleScanError(err, dOptions.PassOnTimeOut)
	}

	if result.LongQueue {
		return dScanner.handleLongQueue(dOptions, result, resolutionErr)
	}

	if dScanner.reportResult(dOptions, result) {
		return FailPipelineErr
	}

	// A non-fatal resolution failure is deliberately surfaced only here, so that
	// its exit code never costs the user the scan results they asked for.
	return resolutionErr
}

// handleLongQueue reports a scan that is still queued once progress polling
// gives up. Passing on it is opt-in via --pass-on-timeout, which also covers
// service access timeouts.
func (dScanner *DebrickedScanner) handleLongQueue(
	options DebrickedOptions,
	result *upload.UploadResult,
	resolutionErr error,
) error {
	fmt.Println("Progress polling terminated due to long scan times. Please try again later")
	fmt.Printf("For full details, visit: %s\n\n", color.BlueString(result.DetailsUrl))

	if options.PassOnTimeOut {
		return resolutionErr
	}

	return LongQueueErr
}

// reportResult renders a completed scan and reports whether a triggered
// automation rule requires the pipeline to fail.
func (dScanner *DebrickedScanner) reportResult(options DebrickedOptions, result *upload.UploadResult) bool {
	WriteApiReplyToJsonFile(options, result)

	fmt.Printf("\n%d vulnerabilities found\n", result.VulnerabilitiesFound)
	fmt.Println("")
	failPipeline := false
	for _, rule := range result.AutomationRules {
		tui.NewRuleCard(os.Stdout, rule).Render()
		failPipeline = failPipeline || (rule.Triggered && rule.FailPipeline())
	}
	fmt.Printf("For full details, visit: %s\n\n", color.BlueString(result.DetailsUrl))

	return failPipeline
}

// isFatalResolutionErr reports whether a resolution error should abort the scan
// before anything is uploaded. Resolution reports non-fatal outcomes as a
// CommandError carrying the exit code the CLI should eventually exit with;
// those let the scan run to completion. Anything else stops it.
func isFatalResolutionErr(err error) bool {
	if err == nil {
		return false
	}

	var cmdErr cmderror.CommandError
	if !errors.As(err, &cmdErr) {
		return true
	}

	return cmdErr.Code == 1
}

func (dScanner *DebrickedScanner) scanReportSBOM(options DebrickedOptions, detailsURL string) error {
	if options.SBOM == "" {
		return nil
	}
	reporter := sbom.Reporter{DebClient: *dScanner.client, FileWriter: io.FileWriter{}}
	repositoryID, commitID, err := reporter.ParseDetailsURL(detailsURL)
	if err != nil {

		return err
	}

	return reporter.Order(sbom.OrderArgs{
		Format:          options.SBOM,
		RepositoryID:    repositoryID,
		CommitID:        commitID,
		Branch:          options.BranchName,
		Vulnerabilities: true,
		Licenses:        true,
		Output:          options.SBOMOutput,
	})
}

func (dScanner *DebrickedScanner) scanResolve(options DebrickedOptions) error {
	resolveOptions := resolution.DebrickedOptions{
		Path:                 options.Path,
		Verbose:              options.Verbose,
		Regenerate:           options.Regenerate,
		Exclusions:           options.Exclusions,
		Inclusions:           options.Inclusions,
		NpmPreferred:         options.NpmPreferred,
		ResolutionStrictness: options.ResolutionStrictness,
	}
	if options.Resolve {
		_, resErr := dScanner.resolver.Resolve([]string{options.Path}, resolveOptions)
		if resErr != nil {
			return resErr
		}
	}

	return nil
}

func (dScanner *DebrickedScanner) scanFingerprint(options DebrickedOptions) error {
	if options.Fingerprint {
		if !(*dScanner.client).IsEnterpriseCustomer(false) {

			return nil
		}
		fingerprints, err := dScanner.fingerprint.FingerprintFiles(
			fingerprint.DebrickedOptions{
				Path:                         options.Path,
				Exclusions:                   append(options.Exclusions, fingerprint.DefaultExclusionsFingerprint()...),
				Inclusions:                   append(options.Inclusions, fingerprint.DefaultInclusionsFingerprint()...),
				MinFingerprintContentLength:  options.MinFingerprintContentLength,
				FingerprintCompressedContent: false,
				Regenerate:                   options.Regenerate > 0,
			},
		)
		if err != nil {
			return err
		}
		err = fingerprints.ToFile(fingerprint.OutputFileNameFingerprints)

		return err
	}

	return nil
}

func (dScanner *DebrickedScanner) scan(options DebrickedOptions, gitMetaObject git.MetaObject) (*upload.UploadResult, error) {

	debug.Log("Running scanFingerprint...", options.Debug)
	err := dScanner.scanFingerprint(options)
	if err != nil {
		return nil, err
	}

	if options.CallGraph {
		debug.Log("Running callgraph generation...", options.Debug)
		javaConfigKwargs := map[string]string{"pm": "maven"}
		javaConfigKwargs["verbose"] = strconv.FormatBool(options.Verbose)
		javaEngine := options.JavaCallgraphEngine
		if javaEngine == "" {
			javaEngine = "soot"
		}
		javaConfigKwargs["java-callgraph-engine"] = javaEngine
		configs := []config.IConfig{
			config.NewConfig("java", []string{}, javaConfigKwargs, true, "maven", options.Version),
			config.NewConfig("golang", []string{}, map[string]string{"pm": "go"}, true, "go", options.Version),
		}
		timeout := options.CallGraphGenerateTimeout
		path := options.Path
		if path == "" {
			path = "."
		}
		resErr := dScanner.callgraph.GenerateWithTimer(
			callgraph.DebrickedOptions{
				Paths:      []string{path},
				Exclusions: options.Exclusions,
				Inclusions: options.Inclusions,
				Configs:    configs,
				Timeout:    timeout,
			},
		)
		if resErr != nil {
			return nil, resErr
		}
	}

	debug.Log("Matching groups...", options.Debug)
	fileGroups, err := dScanner.finder.GetGroups(
		file.DebrickedOptions{
			RootPath:     options.Path,
			Exclusions:   options.Exclusions,
			Inclusions:   options.Inclusions,
			LockFileOnly: false,
			Strictness:   file.StrictAll,
		},
	)
	if err != nil {
		return nil, err
	}

	debug.Log("Starting upload...", options.Debug)
	uploaderOptions := upload.DebrickedOptions{
		FileGroups:             fileGroups,
		GitMetaObject:          gitMetaObject,
		IntegrationsName:       options.IntegrationName,
		CallGraphUploadTimeout: options.CallGraphUploadTimeout,
		VersionHint:            options.VersionHint,
		DebrickedConfig:        dScanner.getDebrickedConfig(options.Path, options.Exclusions, options.Inclusions),
		TagCommitAsRelease:     options.TagCommitAsRelease,
		Experimental:           options.Experimental,
		NoResolve:              !options.Resolve,
	}
	result, err := (*dScanner.uploader).Upload(uploaderOptions)
	if err != nil {
		return nil, err
	}
	err = dScanner.scanReportSBOM(
		options,
		result.DetailsUrl,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dScanner *DebrickedScanner) getDebrickedConfig(path string, exclusions []string, inclusions []string) *upload.DebrickedConfig {
	configPath := dScanner.finder.GetConfigPath(path, exclusions, inclusions)
	if configPath == "" {
		return nil
	}

	return upload.GetDebrickedConfig(configPath)
}

func (dScanner *DebrickedScanner) handleScanError(err error, passOnTimeOut bool) error {
	if err == client.NoResErr && passOnTimeOut {
		fmt.Println(err)

		return nil
	}

	return err
}

// SetWorkingDirectory sets working directory in accordance with the path option
func SetWorkingDirectory(d *DebrickedOptions) error {
	absPath, _ := filepath.Abs(d.Path)
	err := os.Chdir(absPath)
	if err != nil {
		return err
	}
	d.Path = ""
	fmt.Printf("Working directory: %s\n", absPath)

	return nil
}

func UpdatedEmptyCommitName(o *DebrickedOptions) {
	if o.GenerateCommitName && o.CommitName == "" {
		debug.Log("No commit name set, generating commit name", o.Debug)
		o.CommitName = GenerateCommitNameTimestamp()
	}
}

func GenerateCommitNameTimestamp() string {
	return fmt.Sprintf("generated-%d", time.Now().Unix())
}

func MapEnvToOptions(o *DebrickedOptions, env env.Env) {
	if len(o.RepositoryName) == 0 {
		o.RepositoryName = env.Repository
	}
	if len(o.CommitName) == 0 {
		o.CommitName = env.Commit
	}
	if len(o.BranchName) == 0 {
		o.BranchName = env.Branch
	}
	if len(o.CommitAuthor) == 0 {
		o.CommitAuthor = env.Author
	}
	if len(o.RepositoryUrl) == 0 {
		o.RepositoryUrl = env.RepositoryUrl
	}
	if o.IntegrationName == "CLI" {
		if len(env.Integration) != 0 {
			o.IntegrationName = env.Integration
		}
	}
	if len(o.Path) == 0 && len(env.Filepath) > 0 {
		o.Path = env.Filepath
	}
}

func WriteApiReplyToJsonFile(options DebrickedOptions, result *upload.UploadResult) {
	if options.JsonFilePath != "" {
		file, _ := json.MarshalIndent(result, "", " ")
		_ = os.WriteFile(options.JsonFilePath, file, 0600)
	}
}
