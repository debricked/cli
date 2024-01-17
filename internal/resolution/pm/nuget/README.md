# Nuget resolution logic

There are two supported files for resolution of nuget lock files:

### packages.config

We need to convert a `packages.config` file to a `.csproj` file. This is to enable the use of the dotnet restore command
that enables Debricked to parse out transitive dependencies. This may add some additional framework dependencies that
will not show up if we only scan the `packages.config` file. This is done in a few steps:

1. Parse `packages.config` file
2. Run `dotnet --version` to get dotnet version
3. Collect unique target frameworks and packages from the file
4. Create `.nuget.debricked.csproj.temp` file with the collected data

With this done we can move on to the next section

### .csproj

1. Run `dotnet restore <file> --use-lock-file --lock-file-path <lock_file>` in order to restore the dependencies and tools of a project (lock file name can be different depend on which manifest file is being resolved)
2. Cleanup temporary csproj file after lock file is created (for `packages.config` case)
