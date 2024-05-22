# Bitbucket Pipelines
- [Default template](bitbucket-pipelines.yml)

# Private dependencies

Some additional configuration may be required if you use private dependencies not hosted on the default registries, depending on package manager.

## Maven

For maven your `.m2/settings.xml` needs to be configured for the specific registry you wish to use, see the [settings documentation](https://maven.apache.org/settings.html) for more details.
