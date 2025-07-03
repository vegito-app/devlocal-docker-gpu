ARG builder_image=europe-west1-docker.pkg.dev/moov-438615/docker-repository-public/moov-438615:builder
FROM ${builder_image} AS build

ARG non_root_user=dev
USER ${non_root_user}

RUN mkdir -p ${HOME}/src/dev

WORKDIR ${HOME}/src

COPY Makefile .
RUN mkdir -p application
COPY application/application.mk ./application/application.mk
COPY application/nodejs.mk ./application/nodejs.mk
COPY application/go.mk ./application/go.mk

# frontend
RUN mkdir -p application/frontend 
COPY application/frontend/package.json application/frontend/
COPY application/frontend/package-lock.json application/frontend/
COPY application/frontend/frontend.mk application/frontend/
COPY application/application.mk application/

RUN make application-frontend-npm-ci

COPY application/frontend/README.md application/frontend/
COPY application/frontend/public application/frontend/public
COPY application/frontend/src application/frontend/src
COPY application/frontend/webpack.server.js application/frontend/

RUN make -j application-frontend-build application-frontend-bundle

# backend
RUN mkdir -p application/backend

COPY application/backend/backend.mk application/backend/
COPY application/backend/go.sum application/backend/
COPY application/backend/go.mod application/backend/
COPY application/images application/images

RUN make go-application/backend-mod-download

COPY application/backend/btc application/backend/btc
COPY application/backend/firebase application/backend/firebase
COPY application/backend/http application/backend/http
COPY application/backend/internal application/backend/internal
COPY application/backend/log application/backend/log
COPY application/backend/track application/backend/track
COPY application/backend/vault application/backend/vault
COPY application/backend/main.go application/backend/

COPY .git .git
RUN \
    git config --global --add safe.directory ${HOME}/src && \
    make application-backend-install

# artefacts

# # x86_64
# FROM scratch AS final-amd64
# COPY --from=build /lib/x86_64-linux-gnu/libstdc++.so.6 /lib/x86_64-linux-gnu/libstdc++.so.6
# COPY --from=build /lib/x86_64-linux-gnu/libstdc++.so.6.0.30 /lib/x86_64-linux-gnu/libstdc++.so.6.0.30
# COPY --from=build /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
# COPY --from=build /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/x86_64-linux-gnu/libgcc_s.so.1
# COPY --from=build /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
# COPY --from=build /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2
# COPY --from=build /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2

# # arm64
# FROM scratch AS final-arm64
# COPY --from=build /lib/aarch64-linux-gnu/libstdc++.so.6 /lib/aarch64-linux-gnu/libstdc++.so.6
# COPY --from=build /lib/aarch64-linux-gnu/libstdc++.so.6.0.30 /lib/aarch64-linux-gnu/libstdc++.so.6.0.30
# COPY --from=build /lib/aarch64-linux-gnu/libm.so.6 /lib/aarch64-linux-gnu/libm.so.6
# COPY --from=build /lib/aarch64-linux-gnu/libgcc_s.so.1 /lib/aarch64-linux-gnu/libgcc_s.so.1
# COPY --from=build /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libc.so.6
# COPY --from=build /lib/ld-linux-aarch64.so.1 /lib/ld-linux-aarch64.so.1

# ARG TARGETARCH

# Choix de l'étape finale en fonction de l'architecture
# FROM final-${TARGETARCH} AS final

FROM scratch
ARG builder_home=/home/dev
# Actions supplémentaires à réaliser après la copie des bibliothèques
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build ${builder_home}/src/application/frontend/build/ /frontend/build
COPY --from=build ${builder_home}/src/application/frontend/public/ /frontend/public
COPY --from=build ${builder_home}/go/bin/backend /backend

ENV FRONTEND_BUILD_DIR=/frontend/build \
    FRONTEND_PUBLIC_DIR=/frontend/public \
    UI_JAVASCRIPT_SOURCE_FILE=/frontend/build/bundle.js

ENTRYPOINT ["/backend"]