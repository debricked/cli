# GitHub Actions
- [Default template using Docker image](debricked.yml)
- [Default template without Docker](debricked-non-docker.yml)

# Private dependencies

Some additional configuration may be required if you use private dependencies not hosted on the default registries, depending on package manager.

## Maven

For maven your `.m2/settings.xml` needs to be configured for the specific registry you wish to use, see the [settings documentation](https://maven.apache.org/settings.html) for more details.
For more information about maven registries with GitHub Actions see: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-apache-maven-registry
