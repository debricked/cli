# Debricked CLI - Docker
## Docker images
### dev
Contains dev related tools, such as the Go compiler, source code and other.
### cli
A tiny Alpine image that contains `debricked` and is ready to analyse your open source.
### scan
`scan` extends `cli` and is ready to be used in CI/CD tools.
Right now the officially supported CI/CD tools are:
- Azure Pipelines
- Bitbucket Pipelines
- Buildkite
- CircleCI
- GitHub Actions
- GitLab CI/CD
#### How to add CI/CD tool
1. Create a new `elif` statement with a unique tool specific environment variable check.
2. Create a script in `scripts/ci/` which specifies all scan specific environment variables.
    - Variables:
        - `DEBRICKED_SCAN_PATH`
        - `DEBRICKED_SCAN_REPOSITORY`
        - `DEBRICKED_SCAN_COMMIT`
        - `DEBRICKED_SCAN_BRANCH`
        - `DEBRICKED_SCAN_AUTHOR`
        - `DEBRICKED_SCAN_REPOSITORY_URL`
        - `DEBRICKED_SCAN_INTEGRATION`
3. Add test in `tests/`
