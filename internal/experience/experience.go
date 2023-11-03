package experience

import (
	"encoding/json"
	"log"
	"os"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/git"
)

var OutputFileNameExperience = "debricked.experience.json"

type IExperience interface {
	CalculateExperience(rootPath string, exclusions []string) (*Experiences, error)
}

type ExperienceCalculator struct {
	finder file.IFinder
}

func NewExperience(finder file.IFinder) *ExperienceCalculator {
	return &ExperienceCalculator{
		finder: finder,
	}
}

func (e *ExperienceCalculator) CalculateExperience(rootPath string, exclusions []string) (*Experiences, error) {
	log.Println("Calculating experience...")

	repo, repoErr := git.FindRepository(rootPath)
	if repoErr != nil {
		return nil, repoErr
	}

	blamer := git.NewBlamer(repo)

	blames, err := blamer.BlamAllFiles()
	if err != nil {
		return nil, err
	}

	log.Println("Blamed files:", len(blames))
	return nil, nil
}

type Experience struct {
	Author string `json:"author"`
	Email  string `json:"email"`
	Count  int    `json:"count"`
	Symbol string `json:"symbol"`
}

type Experiences struct {
	Entries []Experience `json:"experiences"`
}

func (f *Experiences) ToFile(ouputFile string) error {
	file, err := os.Create(ouputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(f)
}
