FROM debian:bookworm

RUN apt-get update && apt-get install -y \
    adb \
    curl \
    htop \
    iftop \
    jq \
    make \
    netcat-openbsd \
    sudo \
    socat \
    tree \
    procps \
    python3 python3-pip python3-venv \
    && rm -rf /var/lib/apt/lists/*

# Install Stripe CLI via .deb from GitHub Releases
RUN curl -L -o /tmp/stripe-cli.deb https://github.com/stripe/stripe-cli/releases/download/v1.27.0/stripe_1.27.0_linux_amd64.deb \
    && apt-get install -y /tmp/stripe-cli.deb \
    && rm /tmp/stripe-cli.deb

RUN python3 -m venv /opt/venv

RUN /opt/venv/bin/pip install --no-cache-dir \
    firebase-admin \
    robotframework \
    robotframework-appiumlibrary \
    robotframework-requests \
    robotframework-seleniumlibrary \
    robotframework-faker

ENV PATH="/opt/venv/bin:$PATH"

ARG non_root_user=tests
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
COPY entrypoint.sh /usr/bin/tests-entrypoint.sh

ENV LOCAL_APPLICATION_TESTS_DIR=${HOME}/application-tests

ENTRYPOINT [ "tests-entrypoint.sh" ]