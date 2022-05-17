install:
	go install ./cmd/debricked

test:
	test-cli && test-docker-scripts && test-docker
test-cli:
	bash scripts/test_cli.sh
test-docker-scripts:
	bash scripts/test_docker_scripts.sh
test-docker:
	bash scripts/test_docker.sh

docker-build-dev:
	docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target scan .
docker-build-cli:
	docker build -f build/docker/Dockerfile -t debricked/cli:latest --target scan .
docker-build-scan:
	docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .
