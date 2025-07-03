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

socat TCP-LISTEN:8088,fork,reuseaddr TCP:application-tests:8088 > /tmp/socat-application-tests-8088.log 2>&1 &
bg_pids+=("$!")

docker compose -f docker-compose.yml up application-tests 2>&1
