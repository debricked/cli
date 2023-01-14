package scan

type IOptions interface{}

type DebrickedOptions struct {
	Path            string
	Exclusions      []string
	RepositoryName  string
	CommitName      string
	BranchName      string
	CommitAuthor    string
	RepositoryUrl   string
	IntegrationName string
}

type DockerOptions struct {
	Image string
}
