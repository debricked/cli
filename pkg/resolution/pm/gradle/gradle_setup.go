package gradle

import (
	"bufio"
	"embed"
	"fmt"
	"log"
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

type ISetupFile interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile() ([]byte, error)
}

type SetupFile struct{}

type GradleSetup struct {
	gradlewMap        map[string]string
	settingsMap       map[string]string
	subProjectMap     map[string]string
	groovyScriptPath  string
	gradlewOsName     string
	settingsFilenames []string
	GradleProjects    []GradleProject
}

type GradleProject struct {
	dir     string
	gradlew string
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

	return &GradleSetup{
		gradlewMap:        gradlewMap,
		settingsMap:       settingsMap,
		subProjectMap:     subProjectMap,
		groovyScriptPath:  initScript,
		gradlewOsName:     gradlewOsName,
		settingsFilenames: settingsFilenames,
		GradleProjects:    gradleProjects,
	}
}

func (gs *GradleSetup) Setup(files []string, paths []string) {
	writer := writer.FileWriter{}
	err := SetupFile{}.WriteInitFile(gs.groovyScriptPath, writer)
	// gs.setupFilePathMappings(files)
	gs.findGradleProjectFiles(paths)
	// Todo add handling of error
	fmt.Println(err)
	gs.setupGradleProjectMappings()
}

func (gs GradleSetup) findGradleProjectFiles(paths []string) {
	settings := []string{"settings.gradle", "settings.gradle.kts"}
	gradlew := []string{"gradlew"}

	for _, rootPath := range paths {
		err := filepath.Walk(
			rootPath,
			func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !fileInfo.IsDir() {
					for _, setting := range settings {
						if setting == filepath.Base(path) {
							dir, _ := filepath.Abs(filepath.Dir(path))
							file, _ := filepath.Abs(path)
							gs.settingsMap[dir] = file
						}
					}

					for _, gradle := range gradlew {
						if gradle == filepath.Base(path) {
							dir, _ := filepath.Abs(filepath.Dir(path))
							file, _ := filepath.Abs(path)
							gs.gradlewMap[dir] = file
						}
					}
				}

				return nil
			},
		)

		// TODO handle err?
		fmt.Println(err)

	}
}

func (gs *GradleSetup) setupFilePathMappings(files []string) {
	// Setup gradlew / filename mappings (could be better if we could reach the specific inserted paths)
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

func (gs *GradleSetup) setupGradleProjectMappings() {
	// Sort the settingDirs to be in order, hopefully running fewer commands
	settingsDirs := []string{}
	for k, _ := range gs.settingsMap {
		settingsDirs = append(settingsDirs, k)
	}
	sort.Strings(settingsDirs)
	fmt.Println("Sorted settings", settingsDirs)

	// If found settings, run script for finding subprojects on each?
	for _, dir := range settingsDirs {
		// Continue if dir is already subproject of a project
		if _, ok := gs.subProjectMap[dir]; ok {
			continue
		}

		// Setup gradlew, use gradle as default if no gradlew can be found
		gradlew := gs.GetGradleW(dir)
		gradleProject := GradleProject{dir: dir, gradlew: gradlew}

		// Setup subProjectPaths
		gs.setupSubProjectPaths(gradleProject)

		gs.GradleProjects = append(gs.GradleProjects, gradleProject)
	}
}

func (gs *GradleSetup) setupSubProjectPaths(gp GradleProject) {
	// RunMakeFindSubGraphCmd
	factory := CmdFactory{}
	dependenciesCmd, _ := factory.MakeFindSubGraphCmd(gp.dir, gp.gradlew, gs.groovyScriptPath)
	_, err := dependenciesCmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	multiProject := filepath.Join(gp.dir, multiProjectFilename)
	file, err := os.Open(multiProject)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(multiProject)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subProjectPath := scanner.Text()
		gs.subProjectMap[subProjectPath] = gp.dir
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func (gs *GradleSetup) GetGradleW(dir string) string {
	// Get gradlew, if not found in current dir, check if any other gradlew had been found relatively to this path.
	// If so, use that gradlew. If not, use gradle instead of project specific gradlew.
	gradlew := initGradle
	val, ok := gs.gradlewMap[dir]
	if ok {
		gradlew = val
	} else {
		for dirPath, gradlePath := range gs.gradlewMap {
			// potential improvement, sort gradlewMap in longest path first"
			rel, err := filepath.Rel(dirPath, dir)
			isRelative := !strings.HasPrefix(rel, "..") && rel != ".."
			if isRelative == true && err == nil {
				gradlew = gradlePath
				break
			}
		}
	}
	return gradlew
}

func (_ SetupFile) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (sf SetupFile) WriteInitFile(targetFileName string, fileWriter writer.FileWriter) error {
	content, err := sf.ReadInitFile()
	if err != nil {
		return err
	}

	lockFile, err := fileWriter.Create(targetFileName)
	if err != nil {
		return err
	}
	defer lockFile.Close()

	err = fileWriter.Write(lockFile, content)
	return err
}
