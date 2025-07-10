#!/bin/bash

set -eu

# List to hold background job PIDs
bg_pids=()

# Function to kill background jobs when script ends
kill_jobs() {
    echo "Killing background jobs"
    for pid in "$${bg_pids[@]}"; do
        kill "$$pid"
        wait "$$pid" 2>/dev/null
    done
}

# Trap to call kill_jobs on script exit
trap kill_jobs EXIT

make local-firebase-emulators-install
make local-firebase-emulators-start &
bg_pids+=("$!")

# Need localproxy to forward some required port that could not be 
# configured to listen on 0.0.0.0 from firebase.json file.
until nc -z localhost 4400; do
    echo Waiting Firebase Emulator Reserved port at http://localhost:4400/ ;
    sleep 1 ;
done
TARGET_PORT=4400 LISTEN_PORT=4401 localproxy &
bg_pids+=("$!") 

until nc -z localhost 4500; do
    echo Waiting Firebase Emulator Reserved port at http://localhost:4500/ ;
    sleep 1 ;
done
TARGET_PORT=4500 LISTEN_PORT=4501 localproxy &
bg_pids+=("$!") 

until nc -z localhost 9199; do
    echo Waiting Firebase Emulator Reserved port 3 at http://localhost:9199/ ;
    sleep 1 ;
done
TARGET_PORT=9199 LISTEN_PORT=39199 localproxy &

until nc -z localhost 9299; do
    echo Waiting Firebase Emulator Reserved port 3 at http://localhost:9299/ ;
    sleep 1 ;
done
TARGET_PORT=9299 LISTEN_PORT=9399 localproxy &

bg_pids+=("$!") 
until nc -z localhost 9150; do
    echo Waiting Firebase Emulator Reserved port 3 at http://localhost:9150/ ;
    sleep 1 ;
done
TARGET_PORT=9150 LISTEN_PORT=39150 localproxy &
bg_pids+=("$!") 

until nc -z localhost 8085; do
    echo Waiting Firebase Emulator Reserved port 3 at http://localhost:9150/ ;
    sleep 1 ;
done

make local-firebase-emulators-pubsub-init
make local-firebase-emulators-pubsub-check

sleep infinity 