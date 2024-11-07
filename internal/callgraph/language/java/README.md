# Manual build process in case automatic build failure

## Introduction

This guide assists you in manually building your project and preparing it for call graph analysis after an automatic 
build process fails. This ensures that you can still generate a call graph for your Java project despite any issues with
the automatic procedures.

## Manual Build Steps

### Verify Java Version

Ensure you're using at least Java 11, which is required by the call graph generation process. Check your current Java 
version with:

```shell
java -version
```

### Manual Build Command

In the project directory, execute the Maven build command to compile your project and generate the necessary `.class` 
files, while skipping tests to expedite the process:

```shell
mvn package -q -DskipTests -e
```

### Copy External Dependencies

After a successful build, copy your project's external dependencies to the `.debrickedTmpFolder` using the following 
command. This prepares your project for call graph generation by ensuring all dependencies are available:

```shell
mvn -q -B dependency:copy-dependencies -DoutputDirectory=./.debrickedTmpFolder -DskipTests -e
```

## Preparing for Call Graph Generation Without Automatic Build

If the build fails and cannot be resolved, or if you prefer to use your pre-built `.class` files:

- Use the `--no-build` flag with the call graph generation command to bypass the build process:

```shell
debricked callgraph --no-build
```

- Ensure all `.class` files and external dependencies are correctly placed as per the manual build steps.

## Excluding Specific `pom.xml` Files

To exclude specific `pom.xml` files or any other files from the call graph generation, use the exclusion flags provided 
by the CLI tool. Here’s how:

### Using Exclusion Flags

When running the `debricked callgraph` command, use the `-e` or `--exclusions` flag to specify patterns for files you 
wish to exclude. For example, to exclude all `pom.xml` files located in any directory:

```shell
debricked callgraph -e "**/pom.xml"
```

Adjust the pattern as necessary to target specific files or directories.

## Additional Tips

- **Common Errors and Solutions**: Refer to the official documentation for common errors, such as "out of memory" 
  issues. Consider adjusting your Java VM options to allocate more memory for the process if necessary.
- **Advanced Exclusions**: Utilize advanced pattern matching (e.g., `{alt1,...}`) to fine-tune exclusions, especially 
  when working with multiple build configurations or problematic dependencies.

## Conclusion

Following these steps should mitigate issues encountered during the automatic build process, allowing for successful manual preparation for call graph generation. For further assistance or unresolved issues, consult the existing documentation or contact support for more detailed guidance.