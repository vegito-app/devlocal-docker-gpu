ARG builder_image=europe-west1-docker.pkg.local/moov-dev-439608/docker-repository-public/moov-dev-439608:builder-latest
FROM ${builder_image} AS builder
FROM debian:bookworm

COPY --from=builder /usr/local/bin/localproxy /usr/local/bin/localproxy

RUN apt-get update && apt-get install -y \
    curl \
    htop \
    jq \
    make \
    netcat-openbsd \
    openjdk-17-jdk \
    sudo \
    && rm -rf /var/lib/apt/lists/*

ARG non_root_user=firebase

RUN useradd -m ${non_root_user} -u 1000 && echo "${non_root_user}:${non_root_user}" | chpasswd && adduser ${non_root_user} sudo \
    && echo "${non_root_user} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${non_root_user} \
    && chmod 0440 /etc/sudoers.d/${non_root_user} \
    \
    && chown -R ${non_root_user}:${non_root_user} ${HOME}

ENV HOME=/home/${non_root_user}
WORKDIR ${HOME}

USER ${non_root_user}

ENV NVM_DIR=${HOME}/nvm

ARG nvm_version=0.40.1
ARG node_version=22.14.0
RUN set -ex; \
    # 
    mkdir -p ${NVM_DIR} \
    # 
    && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash - \
    && . ${NVM_DIR}/nvm.sh \
    && nvm install ${node_version} \
    && nvm alias default ${node_version} \
    && nvm use default \
    # 
    && npm install -g \
    firebase-tools \
    depcheck \
    npm-check-updates \
    npm-check \
    npm \
    && rm -rf ${HOME}/.npm 

ENV NODE_PATH=$NVM_DIR/versions/node/v${node_version}/lib/node_modules
ENV PATH=$NVM_DIR/versions/node/v${node_version}/bin:$PATH

COPY container-refresh-caches.sh /usr/local/bin/firebase-emulators-refresh-caches.sh
COPY entrypoint.sh /usr/local/bin/firebase-emulators-entrypoint.sh
COPY firebase-emulators-config-create-json.sh /usr/local/bin/firebase-emulators-config-create-json.sh

ENTRYPOINT [ "firebase-emulators-entrypoint.sh" ]