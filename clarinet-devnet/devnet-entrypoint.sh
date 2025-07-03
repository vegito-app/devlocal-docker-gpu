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

LOCAL_CLARINET_DEVNET_CACHE=${PWD}/.containers/clarinet-devnet
mkdir -p $LOCAL_CLARINET_DEVNET_CACHE

# local docker rootless cache 
LOCAL_DOCKERD_ROOTLESS_CACHE=${HOME}/.share/docker
mkdir -p $LOCAL_CLARINET_DEVNET_CACHE/dockerd 
mkdir -p ${HOME}/.share/
ln -s $LOCAL_CLARINET_DEVNET_CACHE/dockerd $LOCAL_DOCKERD_ROOTLESS_CACHE

dockerd-entrypoint.sh --dns=8.8.8.8 --dns=8.8.4.4 &
bg_pids+=("$!")

export DOCKER_HOST=unix:///run/user/1000/docker.sock

until docker info >/dev/null 2>&1; do echo waiting dockerd startup ; sleep 1 ; done

docker info

TARGET_PORT=2375 LISTEN_PORT=2376 localproxy &
bg_pids+=("$!")

mkdir -p ${HOME}/.docker/run
ln -s /run/user/1000/docker.sock ${HOME}/.docker/run/docker.sock

# Needed with github Codespaces which can change the workspace mount specified inside docker-compose.
current_workspace=$(dirname $PWD)
if [ "$current_workspace" != "/workspaces" ] ; then
    sudo ln -s $current_workspace /workspaces
fi

exec "$@"
