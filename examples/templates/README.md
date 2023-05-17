# CI/CD templates
`debricked` can easily be integrated in your CI/CD pipelines. 
Here you can find multiple different templates as to how it can be used within your setup.
- [Argo Workflows](Argo)
- [Azure Pipelines](Azure)
- [Bitbucket Pipelines](Bitbucket)
- [BuildKite](BuildKite)
- [CircleCI](CircleCI)
- [GitHub Actions](GitHub)
- [GitLab CI/CD](GitLab)
- [Travis CI](Travis)

## The dependency tree
In order for us to analyze all dependencies in your project, their versions, and relations, files containing the resolved dependency trees have to be created prior to scanning. Those depend on the package manager used. If files containing the whole dependency tree are not uploaded, that can negatively affect speed and accuracy.

**Example 1:** If npm is used in your project you will have a `package.json` file, but in order for us to scan all your dependencies we need either `package-lock.json` or `yarn.lock` as well.

**Example 2:** If Maven is used in your project you will have a `pom.xml` file, but in order for us to resolve all your dependencies we need a second file, as Maven does not offer a lock file system. Instead, Maven dependency:tree plugin can be used to create a file called `.maven.debricked.lock`

## Debricked CLI dependency resolution
In all templates the manifest file resolution is enabled by default. That means Debricked CLI will attempt to resolve manifest files that belong to package managers that does not offer lock file systems.  
For example, if a `pom.xml` is found by Debricked CLI it will attempt to create `.maven.debricked.lock` automatically.
To disable manifest file resolution, add the flag `--no-resolve`.