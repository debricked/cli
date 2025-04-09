# SBT (Scala) Resolution Logic

The resolution of SBT (Scala Build Tool) dependencies works as follows:

1. Parse the `build.sbt` file to identify any modules
2. Run `sbt makePom` in the project directory to generate a POM file
3. Find the generated `.pom` file (typically in `target/scala-<version>/<project>-<version>.pom`)
4. Copy/rename the `.pom` file to `pom.xml` in the same directory as `build.sbt`
5. Use the existing Maven resolver to handle the `pom.xml` file

This approach allows SBT projects to leverage the existing Maven resolution logic after the POM file is generated.

## Requirements

1. SBT must be installed and available in the PATH
2. The SBT project must be configured to support the `makePom` command (most SBT projects support this by default)
3. Maven dependencies must be resolvable as described in the Maven resolution documentation

## Private Dependencies

Similar to Maven projects, SBT projects might use dependencies from repositories other than the default ones. The
SBT `makePom` command will include these repository configurations in the generated POM file, and then the Maven
resolution process will handle them as described in the Maven README.

## Troubleshooting

If you encounter issues with SBT resolution:

1. Verify SBT is installed and accessible in the PATH
2. Try running `sbt makePom` manually in your project directory to check if it works
3. Inspect the generated `.pom` file in `target/scala-*/` to ensure it contains the correct dependencies
4. Check if any repository authentication is required for your dependencies

## Error Messages

Common error messages and their meanings:

- `SBT wasn't found`: The SBT executable isn't installed or isn't in the PATH
- `Failed to generate Maven POM file`: There was an error during the POM generation process
- `SBT configuration file not found`: The build.sbt file couldn't be found or accessed
- `Failed to parse SBT build file`: The build.sbt file contains syntax errors
- `We weren't able to retrieve one or more dependencies or plugins`: Network issues prevented dependency resolution

## Example Command

```shell
debricked resolve /path/to/scala/project
```

This will:

1. Find all `build.sbt` files in the specified path
2. Generate a POM file for each one
3. Resolve dependencies using the Maven resolver
4. Create `maven.debricked.lock` files with the dependency information