install:
	go install ./cmd/debricked

test:
	make test-cli && make test-static && make test-security && make test-docker
test-cli:
	bash scripts/test_cli.sh
test-static:
	bash scripts/test_static.sh
test-security:
	bash scripts/test_gosec.sh
test-docker:
	bash scripts/test_docker.sh cli

docker-build-dev:
	docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target dev .
docker-build-cli:
	docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .
docker-build-scan:
	docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .
