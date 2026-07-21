.PHONY: install
install:
	sh scripts/install.sh

.PHONY: lint
lint:
	bash scripts/lint.sh

.PHONY: test
test:
	bash scripts/test_cli.sh

.PHONY: test-e2e
test-e2e:
	bash scripts/test_e2e.sh $(type)

# Build the SootUp fat JAR from source.
# The built artifact is placed in:
#   build/sootup/SootUpWrapper/target/sootup-wrapper-1.0.0.jar
.PHONY: build-sootup-jar
build-sootup-jar:
	cd build/sootup/SootUpWrapper && mvn clean package -q

# Build the SootUp fat JAR AND copy it into the Go CLI embedded directory so
# that it is picked up by the //go:embed directive in sootup_handler.go.
.PHONY: build-sootup-jar-embed
build-sootup-jar-embed: build-sootup-jar
	cp build/sootup/SootUpWrapper/target/sootup-wrapper-1.0.0.jar \
		internal/callgraph/language/java/embedded/SootUpWrapper.jar
	@echo "SootUpWrapper.jar embedded at internal/callgraph/language/java/embedded/SootUpWrapper.jar"

.PHONY: test-e2e-docker
docker-build-dev:
	docker build -f build/docker/alpine.Dockerfile -t debricked/cli:dev --target dev .

.PHONY: docker-build-cli
docker-build-cli:
	docker build -f build/docker/alpine.Dockerfile -t debricked/cli:latest --target cli .

.PHONY: docker-build-scan
docker-build-scan:
	docker build -f build/docker/alpine.Dockerfile -t debricked/cli:scan --target scan .

.PHONY: docker-build-cli-resolution
docker-build-cli-resolution:
	docker build -f build/docker/alpine.Dockerfile -t debricked/cli:resolution --target resolution .
