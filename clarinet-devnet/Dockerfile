ARG builder_image=europe-west1-docker.pkg.local/moov-dev-439608/docker-repository-public/moov-dev-439608:builder-latest
FROM ${builder_image} AS project-builder
FROM docker:dind-rootless AS docker-dind

FROM rust:1-alpine3.20 AS rust-build
ARG TARGETPLATFORM
ARG clarinet_version
RUN apk add --no-cache \
    wget musl-dev build-base \
    linux-headers 

RUN wget -nv https://github.com/hirosystems/clarinet/archive/refs/tags/v${clarinet_version}.tar.gz && \
    tar -xf v${clarinet_version}.tar.gz && \
    cd clarinet-${clarinet_version}  && \
    cargo clarinet-install

FROM debian:bookworm

COPY --from=project-builder /usr/local/bin/localproxy /usr/local/bin/localproxy
COPY --from=docker-dind /usr/local/bin/dockerd-entrypoint.sh /usr/local/bin/
COPY --from=docker-dind /usr/local/bin/docker-entrypoint.sh /usr/local/bin/
COPY --from=docker-dind /usr/local/bin/modprobe /usr/local/bin/
COPY --from=rust-build /usr/local/cargo/bin/clarinet /usr/local/bin

RUN apt-get update && apt-get install -y \
    bash-completion \
    curl \
    git \
    htop \
    make \
    sudo \
    wget 

RUN clarinet completion bash && mv clarinet.bash /etc/bash_completion.d/clarinet.bash && \
    echo 'source /etc/bash_completion.d/clarinet.bash' >> ~/.bashrc

ARG docker_version

ARG TARGETPLATFORM
RUN \
    set -eu; \
    \
    mkdir -p /run/user && chmod 1777 /run/user; \
    \ 
    case "$TARGETPLATFORM" in \
    "linux/amd64") \
    url="https://download.docker.com/linux/static/stable/x86_64/docker-${docker_version}.tgz"; \
    ;; \
    "linux/arm64") \
    url="https://download.docker.com/linux/static/stable/aarch64/docker-${docker_version}.tgz"; \
    ;; \
    *) echo >&2 "error: unsupported 'docker.tgz' architecture ($TARGETPLATFORM)"; exit 1 ;; \
    esac; \
    wget -O 'docker.tgz' "$url";  \
    \
    tar --extract \
    --file docker.tgz \
    --strip-components 1 \
    --directory /usr/local/bin/ \
    --no-same-owner \
    'docker/docker' \
    'docker/containerd' \
    'docker/dockerd' \
    'docker/runc' \
    'docker/containerd-shim-runc-v2' \
    'docker/ctr' \
    'docker/docker-proxy' \
    'docker/docker-init' \
    ; \
    rm docker.tgz; \
    \
    docker --version; \
    \
    #
    # DOCKER_ROOTLESS
    # 
    \
    case "$TARGETPLATFORM" in \
    "linux/amd64") \
    url="https://download.docker.com/linux/static/test/x86_64/docker-rootless-extras-${docker_version}.tgz"; \
    ;; \
    "linux/arm64") \
    url="https://download.docker.com/linux/static/test/aarch64/docker-rootless-extras-${docker_version}.tgz"; \
    ;; \
    *) echo >&2 "error: unsupported 'rootless.tgz' architecture ($TARGETPLATFORM)"; exit 1 ;; \
    esac; \
    \
    wget -O 'rootless.tgz' "${url}"; \
    \
    tar --extract \
    --file rootless.tgz \
    --strip-components 1 \
    --directory /usr/local/bin/ \
    'docker-rootless-extras/rootlesskit' \
    'docker-rootless-extras/vpnkit' \
    ; \
    rm rootless.tgz; \
    \
    rootlesskit --version; \
    vpnkit --version;

RUN apt-get update && apt-get install -y \
    iproute2 \
    iptables \
    uidmap 

ARG non_root_user=clarinet
ENV HOME=/home/${non_root_user}
WORKDIR ${HOME}

RUN useradd -m ${non_root_user} -u 1000 && echo "${non_root_user}:${non_root_user}" | chpasswd && adduser ${non_root_user} sudo \
    && echo "${non_root_user} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${non_root_user} \
    && chmod 0440 /etc/sudoers.d/${non_root_user} \
    \
    && chown -R ${non_root_user}:${non_root_user} ${HOME}

# Use Bash
RUN ln -sf /usr/bin/bash /bin/sh
USER ${non_root_user}
COPY entrypoint.sh /usr/local/bin/clarinet-entrypoint.sh
ENTRYPOINT [ "clarinet-entrypoint.sh" ]