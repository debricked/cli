# Azure Pipelines
- [Default template](azure-pipelines.yml)

# Private dependencies

Some additional configuration may be required if you use private dependencies not hosted on the default registries, depending on package manager.

## Maven

For maven your `.m2/settings.xml` needs to be configured for the specific registry you wish to use, see the [settings documentation](https://maven.apache.org/settings.html) for more details.
Using Azure Pipelines this can be done by commenting out the authentication step in the template.
More information about this can be found [here](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/reference/maven-authenticate-v0?view=azure-pipelines).
