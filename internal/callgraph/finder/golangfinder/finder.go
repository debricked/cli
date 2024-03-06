package golanfinder

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/file"
)

type GolangFinder struct{}

func (f GolangFinder) FindRoots(files []string) ([]string, error) {
	var mainFiles []string

	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}

			if node.Name.Name == "main" {
				hasMainFunction := false
				for _, decl := range node.Decls {
					if funcDecl, ok := decl.(*ast.FuncDecl); ok {
						if funcDecl.Name.Name == "main" && funcDecl.Recv == nil &&
							funcDecl.Type.Params.List == nil {
							hasMainFunction = true
							break
						}
					}
				}

				if hasMainFunction {
					mainFiles = append(mainFiles, file)
				}
			}
		}
	}

	return mainFiles, nil
}

func (f GolangFinder) FindDependencyDirs(files []string, findJars bool) ([]string, error) {
	// Not needed for golang
	return []string{}, nil
}

func (f GolangFinder) FindFiles(roots []string, exclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			excluded := file.Excluded(exclusions, path)

			if info.IsDir() && excluded {
				return filepath.SkipDir
			}

			if !info.IsDir() && !excluded && filepath.Ext(path) == ".go" {
				files[path] = true
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	fileList := make([]string, len(files))
	i := 0
	for k := range files {
		fileList[i] = k
		i++
	}

	return fileList, err
}
