package gradle

import (
	"os"
	"path/filepath"
)

type IGradleFileHandler interface {
	Find(paths []string) (map[string]string, map[string]string, error)
}

type GradleFileHandler struct {
	filepath IGradleFilePath
}

type IGradleFilePath interface {
	Walk(root string, walkFn filepath.WalkFunc) error
	Base(path string) string
	Abs(path string) (string, error)
	Dir(path string) string
}

type GradleFilePath struct{}

func (fp GradleFilePath) Walk(root string, walkFn filepath.WalkFunc) error {

	return filepath.Walk(root, walkFn)
}

func (fp GradleFilePath) Base(path string) string {

	return filepath.Base(path)
}

func (fp GradleFilePath) Abs(path string) (string, error) {

	return filepath.Abs(path)
}

func (fp GradleFilePath) Dir(path string) string {

	return filepath.Dir(path)
}

func (gfh GradleFileHandler) Find(paths []string) (map[string]string, map[string]string, error) {
	settings := []string{"settings.gradle", "settings.gradle.kts"}
	gradlew := []string{"gradlew"}
	settingsMap := map[string]string{}
	gradlewMap := map[string]string{}
	for _, rootPath := range paths {
		err := gfh.filepath.Walk(
			rootPath,
			func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil {

					return err
				}
				if !fileInfo.IsDir() {
					for _, setting := range settings {
						if setting == gfh.filepath.Base(path) {
							dir, _ := gfh.filepath.Abs(gfh.filepath.Dir(path))
							file, _ := gfh.filepath.Abs(path)
							settingsMap[dir] = file
						}
					}

					for _, gradle := range gradlew {
						if gradle == gfh.filepath.Base(path) {
							dir, _ := gfh.filepath.Abs(gfh.filepath.Dir(path))
							file, _ := gfh.filepath.Abs(path)
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
