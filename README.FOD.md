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

`debricked` is Debricked's command line interface. It brings open source security, compliance and health to your
project via the command prompt. 

This readme is specific for the use case of scanning Open Source with Debricked through [Fortify on Demand](https://www.microfocus.com/en-us/cyberres/application-security/fortify-on-demand). 
If you are interested in the readme for Debricked standalone, it can be found [here](https://github.com/CarlTern/cli/blob/main/README.md).
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
  <a href="https://github.com/debricked/cli/releases/latest">
    <img src="https://img.shields.io/github/v/release/debricked/cli" />
  </a>
  <a href="https://twitter.com/debrickedab">
    <img src="https://img.shields.io/badge/Twitter-00acee?logo=twitter&logoColor=white" />
  </a>
  <a href="https://www.linkedin.com/company/debricked">
    <img src="https://img.shields.io/badge/LinkedIn-0077B5?logo=linkedin&logoColor=white" />
  </a>
<p align="center">
  <img src="/assets/debricked_resolve.png" alt="CLI Screenshot">
  <br />
</p>

## Installation
Check out the [releases](https://github.com/debricked/cli/releases/latest) page. Choose the asset that is applicable for your system.
Below follow some common ways to install the CLI.
### Linux
```sh
curl -L https://github.com/debricked/cli/releases/latest/download/cli_linux_x86_64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Mac
```sh
curl -L https://github.com/debricked/cli/releases/latest/download/cli_macOS_arm64.tar.gz | tar -xz debricked
```
```sh
./debricked
```
### Windows
1. [Download zip](https://github.com/debricked/cli/releases/latest/download/cli_windows_x86_64.tar.gz)
2. Unpack zip
```sh
.\debricked
```
### Docker
```sh
docker pull debricked/cli
```
## Prepare for scanning open source through Fortify on Demand
If you're looking to scan your Open Source dependencies with Debricked through [Fortify on Demand](https://www.microfocus.com/en-us/cyberres/application-security/fortify-on-demand), 
the Debricked CLI makes the preparation of your payload easy through the `debricked resolve` command. 

> Note: Unlike scanning your open source through Debricked standalone, where the `debricked scan` command can be used, initating a scan through FoD is not possible using the Debricked CLI. You should therefore not use "debricked scan" as a user of FoD.

### What is lock file resolution and why is it needed?
Lock file resolution is the process of using the dependencies requested in a manifest file (which most often is restricted to the direct dependencies of the project) to generate a lock file, containing all direct and indirect/transitive dependencies with locked versions, as well as the relations between the dependencies. 

Getting the complete information for all dependencies, with versions and their relations is important to ensure that Debricked can make a complete and accurate analysis of the project. It will also ensure that the generated SBOM is accurate and that the suggestions made for remediating potential issues are correct. 

Many package managers have support for building and maintaining native lock files from manifest files, while others do not. In most of these cases, there are still native commands that can be used to produce the same information.

### How does the command work?
Once you've installed the CLI, you simply use `debricked resolve` to have Debricked generate the needed lock files for scanning, using FoD. The command identifies all eligible files in the current directory/payload and runs the necessary commands to generate the lock files.

Debricked resolves into native lock files where possible, but uses custom Debricked lock formats when needed. To resolve manifest files (such as package.json and build.gradle) into lock files (eg. yarn.lock and the Debricked lock format gradle.debricked.lock), native commands from the package managers are used, such as `yarn install` and `gradle dependencies`. 

It is therefore important that the package managers are installed, with the right versions, wherever you run the `debricked resolve` command. The best way to achieve this is to run it in a development or build environment.

When the resolution is complete, you will see the list of files that were resolved. If the resolution were to fail, descriptive error messages from the respective package manager 
will be shown in the output.

For more information on how resolution works, check out https://portal.debricked.com/debricked-cli-63/high-performance-scan-faster-more-accurate-and-more-secure-dependency-scanning-293.

### CI/CD integration
If you would rather use `debricked` in your CI/CD pipelines, you can check out the [templates](examples/templates/README.md) for inspiration, replacing `scan` with `resolve`.

## Contributing
Thank you for your interest in making Debricked CLI even better! Read more about contributing to the
project [here](CONTRIBUTING.md).

Also, make sure to check out the [Debricked Portal](https://portal.debricked.com/). There, you can share your great ideas with us! 
