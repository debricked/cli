install:
	go install ./cmd/debricked

test:
	test-cli && test-docker-scripts && test-docker
test-cli:
	bash scripts/test_cli.sh
test-static:
	bash scripts/test_static.sh
test-security:
	bash scripts/test_gosec.sh
test-docker-scripts:
	bash scripts/test_docker_scripts.sh
test-docker:
	bash scripts/test_docker.sh

docker-build-dev:
	docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target dev .
docker-build-cli:
	docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .
docker-build-scan:
	docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .
