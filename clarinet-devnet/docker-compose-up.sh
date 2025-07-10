#!/bin/bash

set -eu

# List to hold background job PIDs
bg_pids=()

# Function to kill background jobs when script ends
kill_jobs() {
    echo "Killing background jobs"
    for pid in "${bg_pids[@]}"; do
        kill "$pid"
        wait "$pid" 2>/dev/null
    done
}

# Trap to call kill_jobs on script exit
trap kill_jobs EXIT

socat TCP-LISTEN:2376,fork,reuseaddr TCP:clarinet-devnet:2376 > /tmp/socat-clarinet-devnet-2376.log 2>&1 &
bg_pids+=("$!")

docker_compose=${LOCAL_DOCKER_COMPOSE:-docker compose -f ${LOCAL_DIR}/docker-compose.yml}

${docker_compose} up clarinet-devnet 2>&1
