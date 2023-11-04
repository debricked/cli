package git

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/debricked/cli/internal/file"
	"github.com/go-git/go-git/v5"
)

type IBlamer interface {
	Blame(path string) (*BlameFile, error)
}

type Blamer struct {
	repository *git.Repository
	inclusions []string
}

func NewBlamer(repository *git.Repository) *Blamer {
	return &Blamer{
		repository: repository,
		inclusions: file.InclusionsExperience(),
	}
}

type BlameFiles struct {
	Files []BlameFile
}

func (b *BlameFiles) ToFile(outputFile string) error {

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, blameFile := range b.Files {
		for _, line := range blameFile.Lines {
			_, err := file.WriteString(fmt.Sprintf("%s,%d,%s,%s\n", blameFile.Path, line.LineNumber, line.Author.Name, line.Author.Email))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type BlameFile struct {
	Lines []BlameLine
	Path  string
}

type BlameLine struct {
	Author     Author
	LineNumber int
}

type Author struct {
	Email string
	Name  string
}

// gitBlameFile runs `git blame --line-porcelain` on the given file and parses the output to populate a slice of BlameLine.
func gitBlameFile(filePath string) ([]BlameLine, error) {
	cmd := exec.Command("git", "blame", "--line-porcelain", filePath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&out)
	var blameLines []BlameLine
	var currentBlame BlameLine

	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "author "):
			currentBlame.Author.Name = strings.TrimPrefix(line, "author ")
		case strings.HasPrefix(line, "author-mail "):
			currentBlame.Author.Email = strings.Trim(strings.TrimPrefix(line, "author-mail "), "<>")

		case strings.HasPrefix(line, "filename "):

			// End of the current commit block
			currentBlame.LineNumber = lineNumber
			blameLines = append(blameLines, currentBlame) // Add the populated BlameLine to the slice.
			currentBlame = BlameLine{}                    // Reset for the next block.
			lineNumber += 1
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return blameLines, nil
}

func (b *Blamer) BlamAllFiles() (*BlameFiles, error) {
	files, err := FindAllTrackedFiles(b.repository)
	if err != nil {
		return nil, err
	}

	blameFiles := make([]BlameFile, 0)

	blameFileChan := make(chan BlameFile, len(files))
	errChan := make(chan error, len(files))

	w, err := b.repository.Worktree()
	if err != nil {
		log.Fatalf("Could not get workdir: %v", err)
	}

	root := w.Filesystem.Root()

	var wg sync.WaitGroup
	for _, fileBlame := range files {

		// Add the root path to the file path
		fileBlameAbsPath := filepath.Join(root, fileBlame)

		if !file.Included(b.inclusions, fileBlame) {
			continue
		}

		wg.Add(1)
		go func(fileBlame string) {
			defer wg.Done()

			blameLines, err := gitBlameFile(fileBlameAbsPath)

			if err != nil {
				errChan <- err
				return
			}

			blameFile := BlameFile{
				Lines: blameLines,
				Path:  fileBlame,
			}

			blameFileChan <- blameFile
		}(fileBlame)
	}

	wg.Wait()
	close(blameFileChan)
	close(errChan)

	for bf := range blameFileChan {
		blameFiles = append(blameFiles, bf)
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &BlameFiles{
		Files: blameFiles,
	}, nil

}
