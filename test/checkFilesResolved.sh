#!/usr/bin/env bash
GREEN='\033[0;32m'
RED='\033[0;31m'

# This script checks if all files in the current directory exist
function check_files_exist() {
    isOk=true
    for file in "$@"; do
        if test -f "$file"; then
            echo -e "${GREEN}File $file OK${GREEN}"
        else
            echo -e "${RED}File $file does not exist!${RED}"
            isOk=false
        fi
    done

    if [ "$isOk" = false ]; then
        exit 1
    fi

}

check_files_exist "test/resolve/testdata/npm/yarn.lock" \
                  "test/resolve/testdata/pip/requirements.txt.pip.debricked.lock" \
                  "test/resolve/testdata/nuget/packagesconfig/packages.config.nuget.debricked.lock" \
                  "test/resolve/testdata/nuget/csproj/packages.lock.json" \
                  "test/resolve/testdata/gradle/gradle.debricked.lock" \
                  "test/resolve/testdata/maven/maven.debricked.lock" \
                  "test/resolve/testdata/gomod/gomod.debricked.lock"