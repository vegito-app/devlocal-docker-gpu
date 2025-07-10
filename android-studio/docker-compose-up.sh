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

socat TCP-LISTEN:5901,fork,reuseaddr TCP:android-studio:5900 >> /tmp/socat-android-studio-5900.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:5037,fork,reuseaddr TCP:android-studio:5037 >> /tmp/socat-android-studio-5037.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9100,fork,reuseaddr TCP:android-studio:9100 >> /tmp/socat-android-studio-9100.log 2>&1 &
bg_pids+=("$!")


docker_compose=${LOCAL_DOCKER_COMPOSE:-docker compose -f ${LOCAL_DIR}/docker-compose.yml}

${docker_compose} up android-studio 2>&1
