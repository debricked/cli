FROM golang:1.23-bookworm AS dev
WORKDIR /cli

ARG DEBIAN_FRONTEND=noninteractive

RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN mkdir -p internal/file/embedded && \
    wget -O internal/file/embedded/supported_formats.json https://debricked.com/api/1.0/open/files/supported-formats
RUN go mod download && go mod verify
COPY . .
RUN make install
CMD [ "debricked" ]

FROM debian:bookworm-slim AS cli-base

ARG DEBIAN_FRONTEND=noninteractive

ENV DEBRICKED_TOKEN=""
RUN apt -y update && apt -y upgrade && apt -y install git && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

# Please update resolution step accordingly when changing this
FROM cli-base AS cli
COPY --from=dev /cli/debricked /usr/bin/debricked

FROM cli AS scan
CMD [ "debricked",  "scan" ]

FROM cli-base AS resolution

ARG DEBIAN_FRONTEND=noninteractive

# Copy Go from the dev stage to avoid Debian package issues
COPY --from=dev /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"

RUN echo "deb http://deb.debian.org/debian unstable main" >> /etc/apt/sources.list && \
    echo "Package: *" >> /etc/apt/preferences && \
    echo "Pin: release a=unstable" >> /etc/apt/preferences && \
    echo "Pin-Priority: -2" >> /etc/apt/preferences

# Uncomment below if testing packages are needed
#RUN echo "deb http://deb.debian.org/debian testing-updates main" >> /etc/apt/sources.list && \
#    echo "deb http://deb.debian.org/debian testing main" >> /etc/apt/sources.list && \
#    echo "Package: *" >> /etc/apt/preferences && \
#    echo "Pin: release a=testing" >> /etc/apt/preferences && \
#    echo "Pin-Priority: -3" >> /etc/apt/preferences

RUN apt -y update && apt -y upgrade && apt -y install curl gnupg unzip && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/apt/keyrings

ENV MAVEN_VERSION="3.9.9"
ENV MAVEN_HOME="/usr/lib/mvn"
ENV PATH="$MAVEN_HOME/bin:$PATH"
RUN curl -fsSLO http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

ENV GRADLE_VERSION="8.10.1"
ENV GRADLE_HOME="/usr/lib/gradle"
ENV PATH="$GRADLE_HOME/gradle-$GRADLE_VERSION/bin:$PATH"
RUN curl -fsSLO https://services.gradle.org/distributions/gradle-$GRADLE_VERSION-bin.zip && \
    unzip gradle-$GRADLE_VERSION-bin.zip -d $GRADLE_HOME && \
    rm gradle-$GRADLE_VERSION-bin.zip

ENV NODE_MAJOR="20"
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt -y update && apt -y upgrade && apt -y install nodejs && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
RUN npm install --global npm@latest && \
    npm install --global yarn && \
    npm install --global bower

RUN npm -v && yarn -v && bower -v

# https://learn.microsoft.com/en-us/dotnet/core/install/linux-scripted-manual#scripted-install
# https://learn.microsoft.com/en-us/dotnet/core/install/linux-debian
# Package manager installs are only supported on the x64 architecture. Other architectures, such as Arm, must install .NET by some other means such as with Snap, an installer script, or through a manual binary installation.
ENV DOTNET_ROOT="/usr/lib/dotnet"
ENV DOTNET_MAJOR="8.0"
ENV PATH="$DOTNET_ROOT:$PATH"
RUN apt -y update && apt -y install libicu72 && \
    apt -y clean && rm -rf /var/lib/apt/lists/*
RUN curl -fsSLO https://dot.net/v1/dotnet-install.sh \
    && chmod u+x ./dotnet-install.sh \
    && ./dotnet-install.sh --channel $DOTNET_MAJOR --install-dir $DOTNET_ROOT \
    && rm ./dotnet-install.sh \
    && dotnet help

RUN apt -y update && apt -y upgrade && apt -y install ca-certificates && \
    apt -y install -t unstable \
    python3.12 \
    python3.12-venv \
    openjdk-21-jdk && \
    apt -y clean && rm -rf /var/lib/apt/lists/* && \
    ln -s /usr/bin/python3.12 /usr/bin/python

RUN dotnet --version && go version

RUN apt update -y && \
    apt install -t unstable lsb-release apt-transport-https ca-certificates software-properties-common -y && \
    curl -o /etc/apt/trusted.gpg.d/php.gpg https://packages.sury.org/php/apt.gpg && \
    sh -c 'echo "deb https://packages.sury.org/php/ bookworm main" > /etc/apt/sources.list.d/php.list' && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

# Add SBT, used for Scala resolution
ENV SBT_VERSION="1.10.11"
ENV SBT_HOME="/usr/lib/sbt"
ENV PATH="$SBT_HOME/bin:$PATH"
RUN curl -fsSLO https://github.com/sbt/sbt/releases/download/v${SBT_VERSION}/sbt-${SBT_VERSION}.tgz && \
    mkdir -p $SBT_HOME && \
    tar -zxvf sbt-${SBT_VERSION}.tgz -C $SBT_HOME --strip-components=1 && \
    rm sbt-${SBT_VERSION}.tgz && \
    ln -s $SBT_HOME/bin/sbt /usr/bin/sbt

RUN sbt --version

RUN apt -y update && apt -y install \
    php8.3 \
    php8.3-curl \
    php8.3-mbstring \
    php8.3-phar && \
    apt -y clean && rm -rf /var/lib/apt/lists/*

RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/bin --filename=composer

RUN ln -sf /usr/bin/python3.12 /usr/bin/python3 && php -v && composer --version && python3 --version

CMD [ "debricked",  "scan" ]

# Put copy at the end to speedup Docker build by caching previous RUNs and run those concurrently
COPY --from=dev /cli/debricked /usr/bin/debricked
