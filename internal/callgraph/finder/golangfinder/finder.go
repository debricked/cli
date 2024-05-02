package golangfinder

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
			isMain, err := f.isMainFile(file)
			if err != nil {
				return nil, err
			}

			if isMain {
				mainFiles = append(mainFiles, file)
			}
		}
	}

	return mainFiles, nil
}

func (f GolangFinder) isMainFile(file string) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}

	if node.Name.Name != "main" {
		return false, nil
	}

	return f.hasMainFunction(node), nil
}

func (f GolangFinder) hasMainFunction(node *ast.File) bool {
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if f.isMainFunction(funcDecl) {
				return true
			}
		}
	}

	return false
}

func (f GolangFinder) isMainFunction(funcDecl *ast.FuncDecl) bool {
	return funcDecl.Name.Name == "main" && funcDecl.Recv == nil && funcDecl.Type.Params.List == nil
}

// Not needed for golang
func (f GolangFinder) FindDependencyDirs(files []string, findJars bool) ([]string, error) {
	return []string{}, nil
}

func (f GolangFinder) FindFiles(roots []string, exclusions []string, inclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			excluded := file.Excluded(exclusions, inclusions, path)

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
