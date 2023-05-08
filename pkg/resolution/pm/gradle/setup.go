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
	initGradle               = "gradle"
	multiProjectFilename     = ".debricked.multiprojects.txt"
	gradleInitScriptFileName = ".gradle-init-script.debricked.groovy"
)

//go:embed gradle-init/gradle-init-script.groovy
var gradleInitScript embed.FS

type ISetup interface {
	Configure(files []string, paths []string) (Setup, error)
}

type Setup struct {
	gradlewMap        map[string]string
	settingsMap       map[string]string
	subProjectMap     map[string]string
	groovyScriptPath  string
	gradlewOsName     string
	settingsFilenames []string
	GradleProjects    []Project
	CmdFactory        ICmdFactory
	MetaFileFinder    IMetaFileFinder
	InitScriptHandler IInitScriptHandler
	Writer            writer.IFileWriter
}

func NewGradleSetup() *Setup {
	groovyScriptPath, _ := filepath.Abs(gradleInitScriptFileName)
	gradlewOsName := "gradlew"
	if runtime.GOOS == "windows" {
		gradlewOsName = "gradlew.bat"
	}

	return &Setup{
		gradlewMap:        map[string]string{},
		settingsMap:       map[string]string{},
		subProjectMap:     map[string]string{},
		groovyScriptPath:  groovyScriptPath,
		gradlewOsName:     gradlewOsName,
		settingsFilenames: []string{"settings.gradle", "settings.gradle.kts"},
		GradleProjects:    []Project{},
		CmdFactory:        CmdFactory{},
		MetaFileFinder:    MetaFileFinder{filepath: FilePath{}},
		InitScriptHandler: InitScriptHandler{},
		Writer:            writer.FileWriter{},
	}
}

func (gs *Setup) Configure(_ []string, paths []string) (Setup, error) {
	err := gs.InitScriptHandler.WriteInitFile(gs.groovyScriptPath, gs.Writer)
	if err != nil {

		return *gs, err
	}
	settingsMap, gradlewMap, err := gs.MetaFileFinder.Find(paths)
	gs.gradlewMap = gradlewMap
	gs.settingsMap = settingsMap
	if err != nil {

		return *gs, err
	}
	err = gs.setupGradleProjectMappings()
	if err != nil && len(err.Error()) > 0 {
		return *gs, err
	}

	return *gs, nil
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

func (gs *Setup) setupSubProjectPaths(gp Project) error {
	dependenciesCmd, _ := gs.CmdFactory.MakeFindSubGraphCmd(gp.dir, gp.gradlew, gs.groovyScriptPath)
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
