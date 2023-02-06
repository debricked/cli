install:
	bash scripts/install.sh

lint:
	bash scripts/test_lint.sh

test:
	bash scripts/test_cli.sh
test-docker:
	bash scripts/test_docker.sh cli

docker-build-dev:
	docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target dev .
docker-build-cli:
	docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .
docker-build-scan:
	docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .
