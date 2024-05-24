# Maven resolution logic

The way resolution of maven lock files works is as follows:

1. Parse `pom.xml` file 
2. Run `mvn dependency:tree -DoutputFile=maven.debricked.lock -DoutputType=tgf --fail-at-end` in order to install all dependencies

The result of the second command above is then written to `maven.debricked.lock` file.

## Private dependencies / Third party repositories

Many maven projects use repositories other than the default central repository, this can be configured in the projects pom.xml.
However, if the repository is not public it could require authentication and some configuration may be needed for the Debricked CLI to be able to resolve dependencies.

In general, the authentication is handled in a file called `settings.xml` in the `.m2` folder (see [settings documentation](https://maven.apache.org/settings.html) for more information). On your build server and locally on development machines this is probably
already setup, but that may not be the case for the environment where the Debricked scan is running, meaning it will fail to resolve.
To fix this the settings file can be manually changed to add the configuration for the required private repositories, or it can be configured by the pipeline provider (such as AWS or Azure).
