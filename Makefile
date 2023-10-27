.PHONY: install
install:
	bash scripts/install.sh

.PHONY: lint
lint:
	bash scripts/lint.sh

.PHONY: test
test:
	bash scripts/test_cli.sh

.PHONY: test-e2e
test-e2e:
	bash scripts/test_e2e.sh $(type)

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
