#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

# Bash history
cat <<'EOF' >> ~/.bashrc
export HISTSIZE=50000
export HISTFILESIZE=100000
EOF

local_container_cache=${LOCAL_ANDROID_STUDIO_CONTAINER_CACHE:-${LOCAL_DIR:-${PWD}}/.containers/dev}
mkdir -p $local_container_cache

# GO
GOPATH=${HOME}/go
rm -rf $GOPATH
mkdir -p ${local_container_cache}/gopath
ln -sf ${local_container_cache}/gopath $GOPATH
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
mkdir -p ${local_container_cache}/npm
ln -sf ${local_container_cache}/npm $NPM_DIR
 
# GCP
GCLOUD_CONFIG=${HOME}/.config/gcloud
mkdir -p $GCLOUD_CONFIG ${local_container_cache}/gcloud
rm -rf $GCLOUD_CONFIG
ln -sf ${local_container_cache}/gcloud $GCLOUD_CONFIG

# Dev container docker socket locally forwarded from distant dockerd host VM
socat TCP-LISTEN:2375,fork UNIX-CONNECT:/var/run/docker.sock &

# Needed with github Codespaces which can change the workspace mount specified inside docker-compose.
current_workspace=$PWD
if [ "$current_workspace" != "$HOST_PWD" ] ; then
    sudo ln -s $current_workspace $HOST_PWD 2>&1 || true
    echo "Linked current workspace $current_workspace to $HOST_PWD"
fi

exec "$@"