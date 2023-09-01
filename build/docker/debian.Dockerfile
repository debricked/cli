FROM golang:1.20-bullseye AS dev
WORKDIR /cli
RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o debricked ./cmd/debricked
ENTRYPOINT ["debricked"]

FROM debian:bullseye-slim AS cli
ENV DEBRICKED_TOKEN=""
RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
ENTRYPOINT [ "debricked",  "scan" ]

FROM cli AS resolution
RUN echo "deb http://ftp.us.debian.org/debian testing-updates main" >> /etc/apt/sources.list && \
    echo "deb http://ftp.us.debian.org/debian testing main" >> /etc/apt/sources.list && \
    echo "Package: *" >> /etc/apt/preferences && \
    echo "Pin: release a=testing" >> /etc/apt/preferences && \
    echo "Pin-Priority: -2" >> /etc/apt/preferences

RUN apt -y update && apt -y upgrade && apt -y install openjdk-11-jre \
    wget \
    unzip \
    python3 \
    python3-scipy \
    ca-certificates \
    curl \
    gnupg \
    python3-pip && \
    apt -y install -t testing golang-1.20 && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/apt/keyrings

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


ENV NODE_MAJOR 18
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt -y update && apt -y upgrade && apt -y install nodejs && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
RUN npm install -g npm@latest
RUN npm install --global yarn
