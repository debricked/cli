package gradle

import (
	"bufio"
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	initGradle           = "gradle"
	multiProjectFilename = ".debricked.multiprojects.txt"
)

//go:embed gradle-init/gradle-init-script.groovy
var gradleInitScript embed.FS

type IGradleSetup interface {
	Setup(files []string, paths []string) (GradleSetup, error)
}

type GradleSetup struct {
	gradlewMap        map[string]string
	settingsMap       map[string]string
	subProjectMap     map[string]string
	groovyScriptPath  string
	gradlewOsName     string
	settingsFilenames []string
	GradleProjects    []GradleProject
	CmdFactory        ICmdFactory
	FileHandler       IGradleFileHandler
	InitFileHandler   IInitFileHandler
	Writer            writer.IFileWriter
}

type GradleProject struct {
	dir     string
	gradlew string
}

type GradleSetupScriptError struct {
	message string
}

type GradleSetupWalkError struct {
	message string
}

type GradleSetupSubprojectError struct {
	message string
}

func (e GradleSetupScriptError) Error() string {

	return e.message
}

func (e GradleSetupWalkError) Error() string {

	return e.message
}

func (e GradleSetupSubprojectError) Error() string {

	return e.message
}

type GradleSetupError []error

func (e GradleSetupError) Error() string {
	var s string
	for _, err := range e {
		s += err.Error() + "\n"
	}

	return s
}

func NewGradleSetup() *GradleSetup {
	initScript, _ := filepath.Abs(".gradle-init-script.debricked.groovy")
	gradlewOsName := "gradlew"
	if runtime.GOOS == "windows" {
		gradlewOsName = "gradlew.bat"
	}
	settingsFilenames := []string{"settings.gradle", "settings.gradle.kts"}
	gradlewMap := map[string]string{}
	settingsMap := map[string]string{}
	subProjectMap := map[string]string{}
	gradleProjects := []GradleProject{}
	CmdFactory := CmdFactory{}
	FileHandler := GradleFileHandler{filepath: GradleFilePath{}}
	InitFileHandler := InitFileHandler{}
	Writer := writer.FileWriter{}

	return &GradleSetup{
		gradlewMap:        gradlewMap,
		settingsMap:       settingsMap,
		subProjectMap:     subProjectMap,
		groovyScriptPath:  initScript,
		gradlewOsName:     gradlewOsName,
		settingsFilenames: settingsFilenames,
		GradleProjects:    gradleProjects,
		CmdFactory:        CmdFactory,
		FileHandler:       FileHandler,
		InitFileHandler:   InitFileHandler,
		Writer:            Writer,
	}
}

func (gs GradleSetup) Setup(files []string, paths []string) (GradleSetup, error) {
	err := gs.InitFileHandler.WriteInitFile(gs.groovyScriptPath, gs.Writer)
	if err != nil {

		return gs, err
	}
	settingsMap, gradlewMap, err := gs.FileHandler.Find(paths)
	gs.gradlewMap = gradlewMap
	gs.settingsMap = settingsMap
	if err != nil {

		return gs, err
	}
	err = gs.setupGradleProjectMappings()
	if err != nil {

		return gs, err
	}

	return gs, nil
}

func (gs *GradleSetup) setupFilePathMappings(files []string) {
	for _, file := range files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		possibleGradlew := filepath.Join(dir, gs.gradlewOsName)
		_, err := os.Stat(possibleGradlew)
		if err == nil {
			gs.gradlewMap[dir] = possibleGradlew
		}
		for _, settingsFilename := range gs.settingsFilenames {
			possibleSettings := filepath.Join(dir, settingsFilename)
			_, err := os.Stat(possibleSettings)
			if err == nil {
				gs.settingsMap[dir] = possibleSettings
			}
		}
	}
}

func (gs *GradleSetup) setupGradleProjectMappings() error {
	var errors GradleSetupError
	settingsDirs := []string{}
	for k := range gs.settingsMap {
		settingsDirs = append(settingsDirs, k)
	}
	sort.Strings(settingsDirs)
	for _, dir := range settingsDirs {
		if _, ok := gs.subProjectMap[dir]; ok {
			continue
		}
		gradlew := gs.GetGradleW(dir)
		gradleProject := GradleProject{dir: dir, gradlew: gradlew}
		err := gs.setupSubProjectPaths(gradleProject)

		if err != nil {
			errors = append(errors, err)
		}
		gs.GradleProjects = append(gs.GradleProjects, gradleProject)
	}

	return GradleSetupSubprojectError{message: errors.Error()}
}

func (gs *GradleSetup) setupSubProjectPaths(gp GradleProject) error {
	dependenciesCmd, _ := gs.CmdFactory.MakeFindSubGraphCmd(gp.dir, gp.gradlew, gs.groovyScriptPath)
	var stderr bytes.Buffer
	dependenciesCmd.Stderr = &stderr
	_, err := dependenciesCmd.Output()
	dependenciesCmd.Stderr = os.Stderr
	if err != nil {
		errorOutput := stderr.String()

		return GradleSetupSubprojectError{message: errorOutput + err.Error()}
	}
	multiProject := filepath.Join(gp.dir, multiProjectFilename)
	file, err := os.Open(multiProject)
	if err != nil {

		return GradleSetupSubprojectError{message: err.Error()}
	}
	defer file.Close()
	defer os.Remove(multiProject)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subProjectPath := scanner.Text()
		gs.subProjectMap[subProjectPath] = gp.dir
	}

	if err := scanner.Err(); err != nil {

		return GradleSetupSubprojectError{message: err.Error()}
	}

	return nil
}

func (gs *GradleSetup) GetGradleW(dir string) string {
	gradlew := initGradle
	val, ok := gs.gradlewMap[dir]
	if ok {
		gradlew = val
	} else {
		for dirPath, gradlePath := range gs.gradlewMap {
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
