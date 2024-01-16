# Gradle resolution logic

The way resolution of gradle lock files works is as follows:

1. Generate init script file for project and subprojects
2. Run `gradle --init-script gradle-init-script.groovy debrickedAllDeps` in order to create dependencies graph 
3. In case permission to execute gradlew is not granted, fallback to PATHs gradle installation is used: `gradle --init-script gradle-init-script.groovy debrickedFindSubProjectPaths` 

The results of the executed command above is then being written into the lock file.
