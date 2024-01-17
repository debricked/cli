# Maven resolution logic

The way resolution of maven lock files works is as follows:

1. Parse `pom.xml` file 
2. Run `mvn dependency:tree -DoutputFile=maven.debricked.lock -DoutputType=tgf --fail-at-end` in order to install all dependencies

The result of the second command above is then written to `maven.debricked.lock` file.
