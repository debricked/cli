FROM golang:1.20-alpine AS dev
WORKDIR /cli
RUN apk update \
  && apk --no-cache --update add git build-base
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o debricked ./cmd/debricked
ENTRYPOINT ["debricked"]

FROM alpine:latest AS cli
ENV DEBRICKED_TOKEN=""
RUN apk add --no-cache git
WORKDIR /root/
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
ENTRYPOINT [ "debricked",  "scan" ]

FROM cli AS resolution
RUN apk --no-cache --update add \
  openjdk11-jre \
  python3 \
  py3-scipy \
  py3-pip \
  go~=1.20 \
  nodejs \ 
  yarn \
  dotnet7-sdk

ENV MAVEN_VERSION 3.9.2
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

ENV GRADLE_VERSION 8.1.1
ENV GRADLE_HOME /usr/lib/gradle
ENV PATH $GRADLE_HOME/gradle-$GRADLE_VERSION/bin:$PATH
RUN wget https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-bin.zip && \
  unzip gradle-$GRADLE_VERSION-bin.zip -d $GRADLE_HOME && \
  rm gradle-$GRADLE_VERSION-bin.zip


# ENV BIN_DIRECTORY /usr/bin
# ENV DOTNET_DIRECTORY /usr/bin/dotnet

# #Install dotnet
# RUN apk add --no-cache --virtual build-dependencies curl \
#   && curl -SL --output dotnet.tar.gz https://download.visualstudio.microsoft.com/download/pr/f8834fef-d2ab-4cf6-abc3-d8d79cfcde11/0ee05ef4af5fe324ce2977021bf9f340/dotnet-sdk-3.1.426-linux-musl-x64.tar.gz \
#   && mkdir -p $DOTNET_DIRECTORY \
#   && tar zxf dotnet.tar.gz -C $DOTNET_DIRECTORY \
#   && chmod +x $DOTNET_DIRECTORY/dotnet \
#   && rm dotnet.tar.gz \
#   && apk del build-dependencies \
#   && rm -r $DOTNET_DIRECTORY/packs $DOTNET_DIRECTORY/sdk/3.1.426/TestHost $DOTNET_DIRECTORY/sdk/3.1.426/Extensions \
#   $DOTNET_DIRECTORY/sdk/3.1.426/FSharp $DOTNET_DIRECTORY/sdk/3.1.426/Roslyn



