FROM golang:1.19-alpine AS dev
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

FROM cli AS scan
ENTRYPOINT [ "debricked",  "scan" ]
