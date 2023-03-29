package gradle

import (
	"os"
	"path/filepath"
)

type IFileFinder interface {
	FindGradleProjectFiles(paths []string) (map[string]string, map[string]string, error)
}

type FileFinder struct {
	filepath IFilePath
}

type IFilePath interface {
	Walk(root string, walkFn filepath.WalkFunc) error
	Base(path string) string
	Abs(path string) (string, error)
	Dir(path string) string
}

type FilePath struct{}

func (fp FilePath) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}

func (fp FilePath) Base(path string) string {
	return filepath.Base(path)
}

func (fp FilePath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (fp FilePath) Dir(path string) string {
	return filepath.Dir(path)
}

func (f FileFinder) FindGradleProjectFiles(paths []string) (map[string]string, map[string]string, error) {
	settings := []string{"settings.gradle", "settings.gradle.kts"}
	gradlew := []string{"gradlew"}
	settingsMap := map[string]string{}
	gradlewMap := map[string]string{}
	for _, rootPath := range paths {
		err := f.filepath.Walk(
			rootPath,
			func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !fileInfo.IsDir() {
					for _, setting := range settings {
						if setting == f.filepath.Base(path) {
							dir, _ := f.filepath.Abs(f.filepath.Dir(path))
							file, _ := f.filepath.Abs(path)
							settingsMap[dir] = file
						}
					}

					for _, gradle := range gradlew {
						if gradle == f.filepath.Base(path) {
							dir, _ := f.filepath.Abs(f.filepath.Dir(path))
							file, _ := f.filepath.Abs(path)
							gradlewMap[dir] = file
						}
					}
				}
				return nil
			},
		)
		if err != nil {
			return nil, nil, GradleSetupWalkError{message: err.Error()}
		}
	}
	return settingsMap, gradlewMap, nil
}
