.PHONY: install
install:
	bash scripts/install.sh

.PHONY: lint
lint:
	bash scripts/lint.sh

.PHONY: test
test:
	bash scripts/test_cli.sh

.PHONY: test-docker
test-docker:
	bash scripts/test_docker.sh cli

.PHONY: test-e2e
test-e2e:
	bash scripts/test_e2e.sh

.PHONY: test-e2e-docker
docker-build-dev:
	docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target dev .

.PHONY: docker-build-cli
docker-build-cli:
	docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .

.PHONY: docker-build-scan
docker-build-scan:
	docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .

.PHONY: docker-build-cli-resolution
docker-build-cli-resolution:
	docker build -f build/docker/Dockerfile -t debricked/cli-resolution:latest --target cli-resolution .
