package gradle

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type Strategy struct {
	files []string
}

func (s Strategy) Invoke() []job.IJob {
	var jobs []job.IJob

	initGradle := "gradle"
	initScript, _ := filepath.Abs(".gradle-init-script.debricked.groovy")

	multiProjectFilename := ".debricked.multiprojects.txt"
	gradlewFilename := "gradlew"
	if runtime.GOOS == "windows" {
		gradlewFilename = "gradlew.bat"
	}
	settingsFilenames := []string{"settings.gradle", "settings.gradle.kts"}

	writer := writer.FileWriter{}
	err := SetupFile{}.WriteInitFile(initScript, writer)
	fmt.Println(err)
	// Todo add handling of error
	gradlewMap := map[string]string{}
	settingsMap := map[string]string{}
	subProjectMap := map[string]string{}

	// Setup gradlew / filename mappings (could be better if we could reach the specific inserted paths)
	for _, file := range s.files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		possibleGradlew := filepath.Join(dir, gradlewFilename)
		_, err := os.Stat(possibleGradlew)
		if err == nil {
			gradlewMap[dir] = possibleGradlew
		}

		for _, settingsFilename := range settingsFilenames {
			possibleSettings := filepath.Join(dir, settingsFilename)
			_, err := os.Stat(possibleSettings)
			if err == nil {
				settingsMap[dir] = possibleSettings
			}
		}
	}

	// Sort the settingDirs to be in order, hopefully running fewer commands
	settingsDirs := []string{}
	for k, _ := range settingsMap {
		settingsDirs = append(settingsDirs, k)
	}
	sort.Strings(settingsDirs)
	fmt.Println("Sorted settings", settingsDirs)

	// If found settings, run script for finding subprojects on each?
	for _, dir := range settingsDirs {
		// dir is already subproject of other project
		if _, ok := subProjectMap[dir]; ok {
			continue
		}

		// Setup gradlew, use gradle as default if no gradlew can be found
		gradlew := initGradle
		val, ok := gradlewMap[dir]
		if ok {
			gradlew = val
		} else {
			for dirPath, _ := range gradlewMap {
				_, err := filepath.Rel(dirPath, dir)
				if err != nil {
					gradlew = val
					break
				}
			}
		}

		factory := CmdFactory{}
		dependenciesCmd, _ := factory.MakeFindSubGraphCmd(dir, gradlew, initScript)
		_, err := dependenciesCmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		multiProject := filepath.Join(dir, multiProjectFilename)
		file, err := os.Open(multiProject)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// defer os.Remove(multiProject)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			subProjectPath := scanner.Text()
			subProjectMap[subProjectPath] = dir
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		jobs = append(jobs, NewJob(dir, gradlew, initScript, factory, writer))
	}

	for _, file := range s.files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		if _, ok := subProjectMap[dir]; ok {
			continue
		}

		gradlew := initGradle
		val, ok := gradlewMap[dir]
		if ok {
			gradlew = val
		}
		factory := CmdFactory{gradlew: gradlew, initScript: initScript}
		jobs = append(jobs, NewJob(dir, factory, writer))
	}

	// Add teardown (remove initfiles ?)
	fmt.Println("SubProjects found", len(subProjectMap))
	fmt.Println("gradlew found", len(gradlewMap))
	fmt.Println("Jobs found", len(jobs))

	return jobs
}

func NewStrategy(files []string) Strategy {
	return Strategy{files}
}
