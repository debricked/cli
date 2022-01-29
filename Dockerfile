FROM golang:1.17-alpine

WORKDIR /debricked

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go install ./debricked.go

ENTRYPOINT ["debricked"]