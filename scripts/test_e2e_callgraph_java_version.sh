if [ -z "$1" ]; then
    DEBRICKED_JAVA_VERSION=11
else
    DEBRICKED_JAVA_VERSION=$1
fi

sed -i "s/<java.version>[0-9]\+<\/java.version>/<java.version>$DEBRICKED_JAVA_VERSION<\/java.version>/" test/callgraph/testdata/mvnproj-build/pom.xml
go test -v ./test/callgraph/maven_test.go
