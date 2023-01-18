package docker

import (
	"errors"
	"github.com/debricked/cli/pkg/ci"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/docker"
	"github.com/debricked/cli/pkg/scan"
	"github.com/debricked/cli/pkg/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDockerCmd(debClient *client.IDebClient) *cobra.Command {
	var u upload.IUploader
	u, _ = upload.NewUploader(debClient)

	var s scan.IScanner
	s = &scan.DockerScanner{
		Uploader:    u,
		SbomFactory: docker.SbomFactory{},
		CiService:   ci.NewService(nil),
	}

	cmd := &cobra.Command{
		Use:   "docker [image]",
		Short: "Scan Docker image",
		Long: `Scan Docker image. [image] must be on any of the following formats:
	docker:yourrepo/yourimage:tag            image from the Docker daemon
	podman:yourrepo/yourimage:tag            image from the Podman daemon
	docker-archive:path/to/yourimage.tar     tarball from disk for archives created from "docker save"
	oci-archive:path/to/yourimage.tar        tarball from disk for OCI archives (from Skopeo or otherwise)
	oci-dir:path/to/yourimage                image directly from a path on disk for OCI layout directories (from Skopeo or otherwise)
	singularity:path/to/yourimage.sif        image directly from a Singularity Image Format (SIF) container on disk
	dir:path/to/yourproject                  image directly from a path on disk (any directory)
	file:path/to/yourproject/file            image directly from a path on disk (any single file)
	registry:yourrepo/yourimage:tag          pull image from a registry (no container runtime required)
`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(s),
		Args: ValidateArgs,
	}

	return cmd
}

func RunE(s scan.IScanner) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		var img string
		if len(args) == 0 {
			img = ""
		} else {
			img = args[0]
		}
		return s.Scan(scan.DockerOptions{Image: img})
	}
}

func ValidateArgs(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("failed to validate argument. Please use one argument")
	}
	return nil
}
