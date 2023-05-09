package finder

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
	dir           string
	gradlew       string
	mainBuildFile string
}

type Setup struct {
	gradlewMap        map[string]string
	settingsMap       map[string]string
	subProjectMap     map[string]string
	groovyScriptPath  string
	gradlewOsName     string
	settingsFilenames []string
	GradleProjects    []Project
	MetaFileFinder    IMetaFileFinder
	Writer            writer.IFileWriter
	InitScriptHandler IInitScriptHandler
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
		gradlewMap:        map[string]string{},
		settingsMap:       map[string]string{},
		subProjectMap:     map[string]string{},
		groovyScriptPath:  groovyScriptPath,
		gradlewOsName:     gradlewOsName,
		settingsFilenames: []string{"settings.gradle", "settings.gradle.kts"},
		GradleProjects:    []Project{},
		MetaFileFinder:    MetaFileFinder{filepath: FilePath{}},
		Writer:            writer,
		InitScriptHandler: ish,
	}
}

func (gs *Setup) Configure(files []string) error {
	err := gs.InitScriptHandler.WriteInitFile()
	if err != nil {

		return err
	}
	settingsMap, gradlewMap, err := gs.MetaFileFinder.Find(files)
	gs.gradlewMap = gradlewMap
	gs.settingsMap = settingsMap
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

func (gs *Setup) setupGradleProjectMappings() error {
	var errors SetupError
	var settingsDirs []string
	for k := range gs.settingsMap {
		settingsDirs = append(settingsDirs, k)
	}
	sort.Strings(settingsDirs)
	for _, dir := range settingsDirs {
		if _, ok := gs.subProjectMap[dir]; ok {
			continue
		}
		gradlew := gs.GetGradleW(dir)
		mainFile := gs.settingsMap[dir]
		gradleProject := Project{dir: dir, gradlew: gradlew, mainBuildFile: mainFile}
		err := gs.setupSubProjectPaths(gradleProject)

		if err != nil {
			errors = append(errors, err)
		}
		gs.GradleProjects = append(gs.GradleProjects, gradleProject)
	}

	return SetupSubprojectError{message: errors.Error()}
}

func MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	path, err := exec.LookPath(gradlew)

	return &exec.Cmd{
		Path: path,
		Args: []string{gradlew, "--init-script", initScript, "debrickedFindSubProjectPaths"},
		Dir:  workingDirectory,
	}, err
}

func (gs *Setup) setupSubProjectPaths(gp Project) error {
	dependenciesCmd, _ := MakeFindSubGraphCmd(gp.dir, gp.gradlew, gs.groovyScriptPath)
	var stderr bytes.Buffer
	dependenciesCmd.Stderr = &stderr
	_, err := dependenciesCmd.Output()
	dependenciesCmd.Stderr = os.Stderr
	if err != nil {
		errorOutput := stderr.String()

		return SetupSubprojectError{message: errorOutput + err.Error()}
	}
	multiProject := filepath.Join(gp.dir, multiProjectFilename)
	file, err := os.Open(multiProject)
	if err != nil {

		return SetupSubprojectError{message: err.Error()}
	}
	defer file.Close()
	defer os.Remove(multiProject)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subProjectPath := scanner.Text()
		gs.subProjectMap[subProjectPath] = gp.dir
	}

	if err := scanner.Err(); err != nil {
		return SetupSubprojectError{message: err.Error()}
	}

	return nil
}

func (gs *Setup) GetGradleW(dir string) string {
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

func FindGradleRoots(files []string) ([]string, error) {
	gradleBuildFiles := FilterFiles(files, "gradle.build(.kts)?")
	gradleSetup := NewGradleSetup()
	err := gradleSetup.Configure(files)
	if err != nil {

		return []string{}, err
	}

	gradleMainDirs := make(map[string]bool)
	for _, gradleProject := range gradleSetup.GradleProjects {
		dir := gradleProject.dir
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
	}
	for _, file := range gradleBuildFiles {
		dir, _ := filepath.Abs(filepath.Dir(file))
		if _, ok := gradleSetup.subProjectMap[dir]; ok {
			continue
		}
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
	}

	roots := []string{}
	for key := range gradleMainDirs {
		roots = append(roots, key)
	}

	return roots, nil
}
