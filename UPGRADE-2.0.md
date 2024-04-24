# Upgrade from 1.X.X to 2.0.0

## Changed behaviours
- Changes default strictness of resolve command to 1 (Exit with code 1 if all files failed to resolve, otherwise exit with code 0 instead of always exiting with code 0)
- File Fingerprint analysis is on by default, rolling roll-out starting with all repositories that start with the letter "C".
- Added inclusion option to commands to force include patterns which are by default ignored by the CLI
- Refactored how exclusion works for fingerprinting to align it with the rest of the CLI, this includes a breaking change for windows where Unix path separators must be used in patterns.

## Runtime upgrades

- Base Docker images have been upgraded from Go 1.21 to 1.22
- In Docker resolution images, the following runtimes have been updated:
   - Upgrade Java from OpenJDK 11 to 21
   - Upgrade Maven from 3.9.2 to 3.9.6
   - Upgrade Gradle from 8.1.1 to 8.7
   - Upgrade Node from 18 to 21
   - Upgrade dotnet from 7.0 to 8.0
   - Upgrade Go from 1.21 to 1.22
   - Upgrade PHP from 8.2 to 8.3
- Debian Docker images have been upgraded from Bullseye (11) to Bookworm (12)
