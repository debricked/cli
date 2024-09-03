<p align="center">
  <p align="center">
    <img width="150" height="150" src="./assets/CLI_logo_1024.png" alt="Logo">
    <h1 align="center"><b>Debricked CLI</b></h1>
    <p align="center">
    Safety through commandline.
      <br />
      <a href="https://debricked.com"><strong>debricked.com »</strong></a>
      <br />
      <br />
    </p>
  </p>
</p>

`debricked` is Debricked's command line interface. It brings open source security, compliance and health to your
project via the command prompt.
<br/>
<br/>
<a href="https://github.com/viktigpetterr/debricked-go-cli/actions/workflows/test.yml">
    <img src="https://github.com/viktigpetterr/debricked-go-cli/actions/workflows/test.yml/badge.svg" />
  </a>
  <a href="https://github.com/viktigpetterr/debricked-go-cli/actions/workflows/debricked.yml">
    <img src="https://github.com/viktigpetterr/debricked-go-cli/actions/workflows/debricked.yml/badge.svg" />
  </a>
    <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
  <a href="https://github.com/debricked/cli/releases/tag/release-v2">
    <img src="https://img.shields.io/github/v/release/debricked/cli" />
  </a>
  <a href="https://twitter.com/debrickedab">
    <img src="https://img.shields.io/badge/Twitter-00acee?logo=twitter&logoColor=white" />
  </a>
  <a href="https://www.linkedin.com/company/debricked">
    <img src="https://img.shields.io/badge/LinkedIn-0077B5?logo=linkedin&logoColor=white" />
  </a>
<p align="center">
  <img src="/assets/cli.png" alt="CLI Screenshot">
  <br />
</p>

## Installation
Check out the [releases](https://github.com/debricked/cli/releases/tag/release-v2) page. Choose the asset that is applicable for your system.
Below follow some common ways to install the CLI.
### Linux
```sh
curl -LsS https://github.com/debricked/cli/releases/download/release-v2/cli_linux_x86_64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Mac
```sh
curl -LsS https://github.com/debricked/cli/releases/download/release-v2/cli_macOS_arm64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Windows
1. [Download zip](https://github.com/debricked/cli/releases/download/release-v2/cli_windows_x86_64.tar.gz)
2. Unpack zip
```sh
.\debricked
```
### Docker
```sh
docker pull debricked/cli:1-resolution-debian
```
## Scan
Once you've installed the CLI, you're ready to scan your project. You can scan a local project, or integrate a scanning mechanism in your CI/CD pipeline.
1. [Sign up to Debricked](https://debricked.com/app/en/register)
2. [Create an access token](https://docs.debricked.com/product/administration/generate-access-token)
3. `debricked scan -t <access-token>`

When the scan is complete, you will see the total number of vulnerabilities found and a list of automation rules that have been evaluated. Read more about automations [here](https://debricked.com/docs/automation/automation-overview.html#automation-overview).

### Docker
To make a scan directly through Docker based on your current working directory, you can use the following command:
```sh
docker run -v $(pwd):/root debricked/cli:1-resolution-debian debricked scan -t <access-token>
```

### CI/CD integration
If you would rather use `debricked` in your CI/CD pipelines, check out the [templates](examples/templates/README.md).

## Contributing
Thank you for your interest in making Debricked CLI even better! Read more about contributing to the
project [here](CONTRIBUTING.md).

Also, make sure to check out the [Debricked Portal](https://portal.debricked.com/). There, you can share your great ideas with us! 

