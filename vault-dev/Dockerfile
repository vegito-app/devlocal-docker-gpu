FROM debian:bookworm

# ARG builder_image=europe-west1-docker.pkg.local/moov-dev-439608/docker-repository-public/moov-dev-439608:builder-latest
# FROM ${builder_image} AS builder

RUN apt-get update && apt-get install -y \
    curl \
    gpg \
    htop \
    jq \
    lsb-release \
    make \
    netcat-openbsd \
    sudo \
    wget

ARG non_root_user=vault

RUN useradd -m ${non_root_user} -u 1000 && echo "${non_root_user}:${non_root_user}" | chpasswd && adduser ${non_root_user} sudo \
    && echo "${non_root_user} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${non_root_user} \
    && chmod 0440 /etc/sudoers.d/${non_root_user} \
    \
    && chown -R ${non_root_user}:${non_root_user} ${HOME}

ENV HOME=/home/${non_root_user}

RUN wget -O - https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    | tee /etc/apt/sources.list.d/hashicorp.list \
    && apt update && apt install -y vault \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${HOME}

USER ${non_root_user}

# For production derivatives of this container, you should add the IPC_LOCK
# capability so that Vault can mlock memory.
COPY entrypoint.sh /usr/local/bin/vault-local-entrypoint.sh
ENTRYPOINT ["vault-local-entrypoint.sh"]
