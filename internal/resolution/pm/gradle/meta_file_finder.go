package gradle

import (
	"os"
	"path/filepath"
)

type IMetaFileFinder interface {
	Find(paths []string) (map[string]string, map[string]string, error)
}

type MetaFileFinder struct {
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

func (finder MetaFileFinder) Find(paths []string) (map[string]string, map[string]string, error) {
	settings := []string{"settings.gradle", "settings.gradle.kts"}
	gradlew := []string{"gradlew"}
	settingsMap := map[string]string{}
	gradlewMap := map[string]string{}
	for _, rootPath := range paths {
		cleanRootPath := filepath.Clean(rootPath)
		err := finder.filepath.Walk(
			cleanRootPath,
			func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil {

					return err
				}
				if !fileInfo.IsDir() {
					for _, setting := range settings {
						if setting == finder.filepath.Base(path) {
							dir, _ := finder.filepath.Abs(finder.filepath.Dir(path))
							file, _ := finder.filepath.Abs(path)
							settingsMap[dir] = file
						}
					}

					for _, gradle := range gradlew {
						if gradle == finder.filepath.Base(path) {
							dir, _ := finder.filepath.Abs(finder.filepath.Dir(path))
							file, _ := finder.filepath.Abs(path)
							gradlewMap[dir] = file
						}
					}
				}

				return nil
			},
		)
		if err != nil {

			return nil, nil, SetupWalkError{message: err.Error()}
		}
	}

	return settingsMap, gradlewMap, nil
}
