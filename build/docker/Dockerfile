FROM golang:1.17-alpine AS dev
WORKDIR /cli
RUN apk add git
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o debricked ./cmd/debricked
ENTRYPOINT ["debricked"]

FROM alpine:latest AS cli
ENV DEBRICKED_TOKEN=""
RUN apk add git
WORKDIR /root/
COPY --from=dev /cli/debricked /usr/bin/debricked
ENTRYPOINT ["debricked"]

FROM cli AS scan
RUN apk add bash
COPY --from=dev /cli/build/docker/scripts /
RUN chmod a+x /entrypoint.sh /setup_env.sh /ci/*
ENTRYPOINT [ "/entrypoint.sh" ]
