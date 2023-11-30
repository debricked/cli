FROM golang:1.20-bullseye AS dev
WORKDIR /cli
RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN mkdir -p internal/file/embedded && \
    wget -O internal/file/embedded/supported_formats.json https://debricked.com/api/1.0/open/files/supported-formats
RUN go mod download && go mod verify
COPY . .
RUN go build -o debricked ./cmd/debricked
ENTRYPOINT ["debricked"]

FROM debian:bullseye-slim AS cli-base
ENV DEBRICKED_TOKEN=""
RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

# Please update resolution step accordingly when changing this
FROM cli-base AS cli
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
ENTRYPOINT [ "debricked",  "scan" ]

FROM cli-base AS resolution
RUN apt -y update && apt -y upgrade && apt -y install curl gnupg unzip && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/apt/keyrings

ENV MAVEN_VERSION 3.9.2
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH
RUN curl -fsSLO http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

ENV GRADLE_VERSION 8.1.1
ENV GRADLE_HOME /usr/lib/gradle
ENV PATH $GRADLE_HOME/gradle-$GRADLE_VERSION/bin:$PATH
RUN curl -fsSLO https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-bin.zip && \
    unzip gradle-$GRADLE_VERSION-bin.zip -d $GRADLE_HOME && \
    rm gradle-$GRADLE_VERSION-bin.zip

ENV NODE_MAJOR 18
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt -y update && apt -y upgrade && apt -y install nodejs && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
RUN npm install --global npm@latest && npm install --global yarn

# https://learn.microsoft.com/en-us/dotnet/core/install/linux-scripted-manual#scripted-install
# https://learn.microsoft.com/en-us/dotnet/core/install/linux-debian
# Package manager installs are only supported on the x64 architecture. Other architectures, such as Arm, must install .NET by some other means such as with Snap, an installer script, or through a manual binary installation.
ENV DOTNET_ROOT /usr/lib/dotnet
ENV DOTNET_MAJOR 7.0
RUN curl -fsSLO https://dot.net/v1/dotnet-install.sh
RUN chmod u+x ./dotnet-install.sh
RUN ./dotnet-install.sh --channel $DOTNET_MAJOR --install-dir $DOTNET_ROOT
RUN rm ./dotnet-install.sh
ENV PATH $DOTNET_ROOT:$PATH

RUN echo "deb http://ftp.us.debian.org/debian testing-updates main" >> /etc/apt/sources.list && \
    echo "deb http://ftp.us.debian.org/debian testing main" >> /etc/apt/sources.list && \
    echo "Package: *" >> /etc/apt/preferences && \
    echo "Pin: release a=testing" >> /etc/apt/preferences && \
    echo "Pin-Priority: -2" >> /etc/apt/preferences

RUN apt -y update && apt -y upgrade && apt -y install openjdk-11-jre \
    python3 \
    python3-venv \
    ca-certificates \
    python3-pip && \
    apt -y install -t testing golang-1.20 && \
    apt -y clean && rm -rf /var/lib/apt/lists/* && \
    # Symlink pip3 to pip, we assume that "pip" works in CLI
    ln -sf /usr/bin/pip3 /usr/bin/pip && \
    ln -sf /usr/bin/python3 /usr/bin/python && \
    # Symlink go binary to bin directory which is in path
    ln -s /usr/lib/go-1.20/bin/go /usr/bin/go

RUN dotnet --version

RUN apt update -y && \
    apt install lsb-release apt-transport-https ca-certificates software-properties-common -y && \
    curl -o /etc/apt/trusted.gpg.d/php.gpg https://packages.sury.org/php/apt.gpg && \
    sh -c 'echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/php.list' && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN apt -y update && apt -y install \
    php8.2 \
    php8.2-curl \
    php8.2-mbstring \
    php8.2-phar && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/bin --filename=composer

RUN php -v && composer --version

# Put copy at the end to speedup Docker build by caching previous RUNs and run those concurrently
COPY --from=dev /cli/debricked /usr/bin/debricked
