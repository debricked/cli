<p align="center">
  <a href="#"/>
  <p align="center">
    <img width="150" height="150" src="/assets/CLI_logo_1024.png" alt="Logo">
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

`debricked` is Debricked's own command line interface. It brings open source security, compliance and health to your
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
  <a href="https://github.com/debricked/cli#beta-software-%EF%B8%8F">
    <img src="https://img.shields.io/badge/stability-beta-33bbff.svg" />
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

## Beta software ⚠️
This product is not in a stable phase. Breaking changes can occur at any time.

## Installation
Check out the [releases](https://github.com/debricked/cli/releases/latest) page. Choose the asset that is applicable for your system.
Bellow follow some common ways to install the CLI.
### Linux
```sh
curl -L https://github.com/debricked/cli/releases/download/v0.0.7/cli_0.0.7_linux_x86_64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Mac
```sh
curl -L https://github.com/debricked/cli/releases/download/v0.0.7/cli_0.0.7_macOS_arm64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Windows
1. [Download zip](https://github.com/debricked/cli/releases/download/v0.0.7/cli_0.0.7_windows_x86_64.tar.gz)
2. Unpack zip
```sh
.\debricked
```
### Docker
```sh
docker pull debricked/cli
```
## Scan
Once you've installed the CLI, you're ready to scan your project. You can scan a local project, or integrate a scanning mechanism in your CI/CD pipeline.
1. [Sign up to Debricked](https://debricked.com/app/en/register)
2. [Create an access token](https://debricked.com/docs/administration/access-tokens.html#creating-access-tokens)
3. `debricked scan -t <access-token>`

If you would rather use Debricked CLI in your CI/CD pipelines, check out the [docs](https://debricked.com/docs/integrations/ci-build-systems/).

When the scan is complete, you will see the total number of vulnerabilities found and a list of automation rules that have been evaluated. Read more about automations [here](https://debricked.com/docs/automation/automation-overview.html#automation-overview).

## Contributing
Thank you for your interest in making Debricked CLI even better! Read more about contributing to the
project [here](CONTRIBUTING.md).

Also, make sure to check out the [Debricked Portal](https://portal.debricked.com/). There, you can share your great ideas with us! 

