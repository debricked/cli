package docker

import (
	"fmt"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/cmd/syft/cli/eventloop"
	"github.com/anchore/syft/cmd/syft/cli/options"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	"log"
	"os"
)

const SbomFileName = "sbom.json"

type ISbomFactory interface {
	Make(image *Image) (*os.File, error)
}

type SbomFactory struct{}

func (sbomFactory SbomFactory) Make(img *Image) (*os.File, error) {
	writer, err := options.MakeWriter([]string{"cyclonedx-json"}, SbomFileName, SbomFileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("unable to write to report destination: %s\n", err.Error())
		}
	}()

	// could be an image or a directory, with or without a scheme
	si, err := source.ParseInputWithName(img.Name(), "", true, fmt.Sprintf("docker/%s", img.Name()))
	if err != nil {
		return nil, fmt.Errorf("could not generate source input for packages command: %w", err)
	}
	registryOptions := &image.RegistryOptions{
		InsecureSkipTLSVerify: false,
		InsecureUseHTTP:       false,
	}
	src, cleanup, err := source.New(*si, registryOptions, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct source from user input %q: %w", si.UserInput, err)
	}
	if cleanup != nil {
		defer cleanup()
	}
	img.Digest = string(src.ID())
	fmt.Printf("Successfully loaded image: %s\n", img.Name())

	s, err := GenerateSBOM(src)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("no SBOM produced for %q", si.UserInput)
	}

	err = writer.Write(*s)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully parsed: %s\n", img.Name())

	return os.OpenFile(SbomFileName, os.O_RDONLY, 0)
}

func GenerateSBOM(src *source.Source) (*sbom.SBOM, error) {
	task, err := generateCatalogPackagesTask()
	if err != nil {
		return nil, err
	}

	s := sbom.SBOM{
		Source: src.Metadata,
		Descriptor: sbom.Descriptor{
			Name:    "debricked",
			Version: "",
		},
	}

	relationships, err := task(&s.Artifacts, src)
	s.Relationships = append(s.Relationships, relationships...)

	return &s, nil
}

func generateCatalogPackagesTask() (eventloop.Task, error) {
	catalogerConfig := cataloger.Config{
		Search: cataloger.SearchConfig{
			IncludeIndexedArchives:   true,
			IncludeUnindexedArchives: false,
			Scope:                    "Squashed",
		},
		Parallelism: 1,
	}
	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		packageCatalog, relationships, theDistro, err := syft.CatalogPackages(src, catalogerConfig)
		if err != nil {
			return nil, err
		}

		results.PackageCatalog = packageCatalog
		results.LinuxDistribution = theDistro

		return relationships, nil
	}

	return task, nil
}
