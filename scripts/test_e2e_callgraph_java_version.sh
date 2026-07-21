if [ -z "$1" ]; then
    DEBRICKED_JAVA_VERSION=11
else
    DEBRICKED_JAVA_VERSION=$1
fi

POM_PATH="test/callgraph/testdata/mvnproj-build/pom.xml"
POM_BAK="${POM_PATH}.bak"

cp "$POM_PATH" "$POM_BAK"

sed "s/<java.version>[0-9]\+<\/java.version>/<java.version>$DEBRICKED_JAVA_VERSION<\/java.version>/" "$POM_PATH" > "${POM_PATH}.tmp"
mv "${POM_PATH}.tmp" "$POM_PATH"

go test -v ./test/callgraph/maven_test.go
test_exit=$?

mv "$POM_BAK" "$POM_PATH"

exit $test_exit
