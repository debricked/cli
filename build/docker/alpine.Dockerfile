FROM golang:1.23-alpine AS dev
WORKDIR /cli
RUN apk --no-cache --update add git build-base
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN mkdir -p internal/file/embedded && \
    wget -O internal/file/embedded/supported_formats.json https://debricked.com/api/1.0/open/files/supported-formats
RUN apk add --no-cache make curl && make install && apk del make curl
CMD [ "debricked" ]

FROM alpine:latest AS cli-base
ENV DEBRICKED_TOKEN=""
RUN apk add --no-cache git
WORKDIR /root/

# Please update resolution step accordingly when changing this
FROM cli-base AS cli
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
CMD [ "debricked",  "scan" ]

FROM cli-base AS resolution
ENV MAVEN_VERSION="3.9.9"
ENV MAVEN_HOME="/usr/lib/mvn"
ENV PATH="$MAVEN_HOME/bin:$PATH"
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

ENV GRADLE_VERSION="8.10.1"
ENV GRADLE_HOME="/usr/lib/gradle"
ENV PATH="$GRADLE_HOME/gradle-$GRADLE_VERSION/bin:$PATH"
RUN wget https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-bin.zip && \
  unzip gradle-$GRADLE_VERSION-bin.zip -d $GRADLE_HOME && \
  rm gradle-$GRADLE_VERSION-bin.zip

# g++ needed to compile python packages with C dependencies (numpy, scipy, etc.)
RUN apk --no-cache --update add \
  openjdk21-jdk \
  python3 \
  py3-scipy \
  py3-pip \
  nodejs \
  npm \
  yarn \
  g++ \
  curl

RUN apk --no-cache --update add dotnet8-sdk go~=1.23 --repository=https://dl-cdn.alpinelinux.org/alpine/v3.20/community

RUN dotnet --version && npm -v && yarn -v

RUN npm install --global bower && bower -v

RUN apk add --no-cache \
    git \
    php83 \
    php83-curl \
    php83-mbstring \
    php83-openssl \
    php83-phar \
    && ln -sf /usr/bin/php83 /usr/bin/php

RUN apk add --no-cache --virtual build-dependencies curl && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/bin --filename=composer \
    && apk del build-dependencies

RUN php -v && composer --version

CMD [ "debricked",  "scan" ]

# Put copy at the end to speedup Docker build by caching previous RUNs and run those concurrently
COPY --from=dev /cli/debricked /usr/bin/debricked
