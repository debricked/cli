package gradle

import (
	"bufio"
	"bytes"
	"embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/debricked/cli/pkg/io/writer"
)

const (
	initGradle               = "gradle"
	multiProjectFilename     = ".debricked.multiprojects.txt"
	gradleInitScriptFileName = ".gradle-init-script.debricked.groovy"
)

//go:embed embeded/gradle-init-script.groovy
var gradleInitScript embed.FS

type ISetup interface {
	Configure(files []string) (Setup, error)
}

type Project struct {
	Dir           string
	Gradlew       string
	MainBuildFile string
}

type Setup struct {
	GradlewMap        map[string]string
	SettingsMap       map[string]string
	SubProjectMap     map[string]string
	GroovyScriptPath  string
	GradlewOsName     string
	SettingsFilenames []string
	GradleProjects    []Project
	MetaFileFinder    IMetaFileFinder
	Writer            writer.IFileWriter
	InitScriptHandler IInitScriptHandler
	CmdFactory        ICmdFactory
}

func NewGradleSetup() *Setup {
	groovyScriptPath, _ := filepath.Abs(gradleInitScriptFileName)
	gradlewOsName := "gradlew"
	if runtime.GOOS == "windows" {
		gradlewOsName = "gradlew.bat"
	}
	writer := writer.FileWriter{}
	ish := InitScriptHandler{groovyScriptPath, "embeded/gradle-init-script.groovy", writer}

	return &Setup{
		GradlewMap:        map[string]string{},
		SettingsMap:       map[string]string{},
		SubProjectMap:     map[string]string{},
		GroovyScriptPath:  groovyScriptPath,
		GradlewOsName:     gradlewOsName,
		SettingsFilenames: []string{"settings.gradle", "settings.gradle.kts"},
		GradleProjects:    []Project{},
		MetaFileFinder:    MetaFileFinder{filepath: FilePath{}},
		Writer:            writer,
		InitScriptHandler: ish,
		CmdFactory:        CmdFactory{},
	}
}

func (gs *Setup) Configure(files []string) error {
	err := gs.InitScriptHandler.WriteInitFile()
	if err != nil {

		return err
	}
	settingsMap, gradlewMap, err := gs.MetaFileFinder.Find(files)
	gs.GradlewMap = gradlewMap
	gs.SettingsMap = settingsMap
	if err != nil {

		return err
	}
	err = gs.setupGradleProjectMappings()
	if err != nil && len(err.Error()) > 0 {

		return err
	}

	return nil
}

func (gs *Setup) setupFilePathMappings(files []string) {
	for _, file := range files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		possibleGradlew := filepath.Join(dir, gs.GradlewOsName)
		_, err := os.Stat(possibleGradlew)
		if err == nil {
			gs.GradlewMap[dir] = possibleGradlew
		}
		for _, settingsFilename := range gs.SettingsFilenames {
			possibleSettings := filepath.Join(dir, settingsFilename)
			_, err := os.Stat(possibleSettings)
			if err == nil {
				gs.SettingsMap[dir] = possibleSettings
			}
		}
	}
}

func (gs *Setup) setupGradleProjectMappings() error {
	var errors SetupError
	var settingsDirs []string
	for k := range gs.SettingsMap {
		settingsDirs = append(settingsDirs, k)
	}
	sort.Strings(settingsDirs)
	for _, dir := range settingsDirs {
		if _, ok := gs.SubProjectMap[dir]; ok {
			continue
		}
		gradlew := gs.GetGradleW(dir)
		mainFile := gs.SettingsMap[dir]
		gradleProject := Project{Dir: dir, Gradlew: gradlew, MainBuildFile: mainFile}
		err := gs.setupSubProjectPaths(gradleProject)

		if err != nil {
			errors = append(errors, err)
		}
		gs.GradleProjects = append(gs.GradleProjects, gradleProject)
	}

	return SetupSubprojectError{message: errors.Error()}
}

type ICmdFactory interface {
	MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error)
}
type CmdFactory struct{}

func (cf CmdFactory) MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	path, err := exec.LookPath(gradlew)

	return &exec.Cmd{
		Path: path,
		Args: []string{gradlew, "--init-script", initScript, "debrickedFindSubProjectPaths"},
		Dir:  workingDirectory,
	}, err
}

func (gs *Setup) setupSubProjectPaths(gp Project) error {
	dependenciesCmd, _ := gs.CmdFactory.MakeFindSubGraphCmd(gp.Dir, gp.Gradlew, gs.GroovyScriptPath)
	var stderr bytes.Buffer
	dependenciesCmd.Stderr = &stderr
	_, err := dependenciesCmd.Output()
	dependenciesCmd.Stderr = os.Stderr
	if err != nil {
		errorOutput := stderr.String()

		return SetupSubprojectError{message: errorOutput + err.Error()}
	}
	multiProject := filepath.Join(gp.Dir, multiProjectFilename)
	file, err := os.Open(multiProject)
	if err != nil {

		return SetupSubprojectError{message: err.Error()}
	}
	defer file.Close()
	defer os.Remove(multiProject)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subProjectPath := scanner.Text()
		gs.SubProjectMap[subProjectPath] = gp.Dir
	}

	if err := scanner.Err(); err != nil {
		return SetupSubprojectError{message: err.Error()}
	}

	return nil
}

func (gs *Setup) GetGradleW(dir string) string {
	gradlew := initGradle
	val, ok := gs.GradlewMap[dir]
	if ok {
		gradlew = val
	} else {
		for dirPath, gradlePath := range gs.GradlewMap {
			// potential improvement, sort gradlewMap in longest path first"
			rel, err := filepath.Rel(dirPath, dir)
			isRelative := !strings.HasPrefix(rel, "..") && rel != ".."
			if isRelative && err == nil {
				gradlew = gradlePath

				break
			}
		}
	}

	return gradlew
}
