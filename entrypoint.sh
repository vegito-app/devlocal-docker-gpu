#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

# Bash history
cat <<'EOF' >> ~/.bashrc
export HISTSIZE=50000
export HISTFILESIZE=100000
EOF

DEV_CONTAINER_CACHE=${PWD}/.containers/dev
mkdir -p $DEV_CONTAINER_CACHE

# GO
GOPATH=${HOME}/go
rm -rf $GOPATH
mkdir -p ${DEV_CONTAINER_CACHE}/gopath
ln -sf ${DEV_CONTAINER_CACHE}/gopath $GOPATH
cat <<'EOF' >> ~/.bashrc
export GOARCH=$(dpkg --print-architecture)
EOF

cat <<'EOF' >> ~/.bashrc
alias hi="docker run --rm -it --privileged -v /:/host --network=host europe-west1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT_ID}/docker-repository-public/${GOOGLE_CLOUD_PROJECT_ID}:builder-latest sudo chroot /host iftop -i eno1"
alias hh="docker run --rm -it --privileged -v /:/host --network=host europe-west1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT_ID}/docker-repository-public/${GOOGLE_CLOUD_PROJECT_ID}:builder-latest sudo chroot /host htop"
alias h="docker run --rm -it --privileged -v /:/host --network=host europe-west1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT_ID}/docker-repository-public/${GOOGLE_CLOUD_PROJECT_ID}:builder-latest sudo chroot /host htop"
alias r="docker run --rm -it --privileged -v /:/host --network=host europe-west1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT_ID}/docker-repository-public/${GOOGLE_CLOUD_PROJECT_ID}:builder-latest sudo chroot /host"
EOF

# NPM
NPM_DIR=${HOME}/.npm
[ -d $NPM_DIR ] && mv $NPM_DIR ${NPM_DIR}_back
mkdir -p ${DEV_CONTAINER_CACHE}/npm
ln -sf ${DEV_CONTAINER_CACHE}/npm $NPM_DIR
 
# GCP
GCLOUD_CONFIG=${HOME}/.config/gcloud
mkdir -p $GCLOUD_CONFIG ${DEV_CONTAINER_CACHE}/gcloud
rm -rf $GCLOUD_CONFIG
ln -sf ${DEV_CONTAINER_CACHE}/gcloud $GCLOUD_CONFIG

# Dev container docker socket locally forwarded from distant dockerd host VM
socat TCP-LISTEN:2375,fork UNIX-CONNECT:/var/run/docker.sock &

# Needed with github Codespaces which can change the workspace mount specified inside docker-compose.
current_workspace=$(dirname $PWD)
if [ "$current_workspace" != "/workspaces" ] ; then
    sudo ln -s $current_workspace /workspaces
fi

exec "$@"