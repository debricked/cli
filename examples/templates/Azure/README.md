# Azure Pipelines
- [Default template](azure-pipelines.yml)

# Private dependencies

Some additional configuration may be required if you use private dependencies not hosted on the default registries, depending on package manager.

## Maven

For maven your `.m2/settings.xml` needs to be configured for the specific registry you wish to use, see the [settings documentation](https://maven.apache.org/settings.html) for more details.
Using Azure Pipelines this can be done by commenting out the authentication step in the template.
More information about this can be found [here](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/reference/maven-authenticate-v0?view=azure-pipelines).

## NuGet

Just like with maven you need to configure access to private registries or sources,
see [NuGet's documentation on source mapping](https://learn.microsoft.com/en-us/nuget/consume-packages/package-source-mapping) for more information 
(and the [authenticated feeds documentation](https://learn.microsoft.com/en-us/nuget/consume-packages/consuming-packages-authenticated-feeds) for details).
When using Azure pipelines your authenticated feeds can be accessed by commenting out the NuGet parts in the default template
(see [nuget authenticate documentation](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/reference/nuget-authenticate-v1?view=azure-pipelines) for more information).
