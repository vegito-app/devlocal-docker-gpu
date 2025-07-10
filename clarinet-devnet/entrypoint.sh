#!/bin/bash

set -eu

# List to hold background job PIDs
bg_pids=()

# Function to kill background jobs when script ends
kill_jobs() {
    echo "Killing clarinet background jobs"
    for pid in "${bg_pids[@]}"; do
        kill "$pid"
        wait "$pid" 2>/dev/null
    done
}

# Trap to call kill_jobs on script exit
trap kill_jobs EXIT

# Local Container Cache
local_container_cache=${LOCAL_CLARINET_DEVNET_CONTAINER_CACHE:-${LOCAL_DIR:-${PWD}}/.containers/clarinet-devnet}
mkdir -p $local_container_cache

# local docker rootless cache 
LOCAL_DOCKERD_ROOTLESS_CACHE=${HOME}/.share/docker
mkdir -p $local_container_cache/dockerd
mkdir -p ${HOME}/.share/
ln -s $local_container_cache/dockerd $LOCAL_DOCKERD_ROOTLESS_CACHE

dockerd-entrypoint.sh --dns=8.8.8.8 --dns=8.8.4.4 &
bg_pids+=("$!")

LOCAL_USER_ID=$(id -u)

export DOCKER_HOST=unix:///run/user/$LOCAL_USER_ID/docker.sock

until docker info >/dev/null 2>&1; do echo waiting dockerd startup ; sleep 1 ; done

docker info

TARGET_PORT=2375 LISTEN_PORT=2376 localproxy &
bg_pids+=("$!")

mkdir -p ${HOME}/.docker/run
ln -s /run/user/$LOCAL_USER_ID/docker.sock ${HOME}/.docker/run/docker.sock


# Bash history
rm -f ~/.bash_history
ln -sfn ${local_container_cache}/bash_history ~/.bash_history
cat <<EOF >> ~/.bashrc
export HISTSIZE=50000
export HISTFILESIZE=100000
export DOCKER_HOST=unix:///run/user/${LOCAL_USER_ID:-1000}/docker.sock
export DOCKER_CONFIG=${local_container_cache}/.docker
export DOCKER_BUILDKIT=1
EOF

exec "$@"

# Set inotify watches limit
echo fs.inotify.max_user_watches=524288 |  sudo tee -a /etc/sysctl.conf; sudo sysctl -p
# Set inotify watches limit for rootless dockerd
echo fs.inotify.max_user_watches=524288 | sudo tee -a /run/user/$LOCAL_USER_ID/sysctl.conf
sudo sysctl -p /run/user/$LOCAL_USER_ID/sysctl.conf