#!/bin/bash

set -eu

# 📌 List of PIDs of background processes
bg_pids=()

# 🧹 Function called at the end of the script to kill background processes
kill_jobs() {
    echo "🧼 Cleaning up background processes..."
    for pid in "${bg_pids[@]}"; do
        kill "$pid" || true
        wait "$pid" 2>/dev/null || true
    done
}

# 🚨 Register cleanup function to run on script exit
trap kill_jobs EXIT

case "${LOCAL_ANDROID_GPU_MODE}" in
    "host")
        display-start-xpra.sh &
        bg_pids+=("$!")
        ;;
    *)
        display-start.sh &
        bg_pids+=("$!")
        ;;
esac

[ "${LOCAL_ANDROID_STUDIO_ENV_SETUP}" = "true" ] && \
    studio-caches-refresh.sh

# Forward firebase-emulators to container as localhost
socat TCP-LISTEN:9299,fork,reuseaddr TCP:firebase-emulators:9399 > /tmp/socat-firebase-emulators-9399.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:4500,fork,reuseaddr TCP:firebase-emulators:4501 > /tmp/socat-firebase-emulators-4501.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:4400,fork,reuseaddr TCP:firebase-emulators:4401 > /tmp/socat-firebase-emulators-4401.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9000,fork,reuseaddr TCP:firebase-emulators:9000 > /tmp/socat-firebase-emulators-9000.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9099,fork,reuseaddr TCP:firebase-emulators:9099 > /tmp/socat-firebase-emulators-9099.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9150,fork,reuseaddr TCP:firebase-emulators:9150 > /tmp/socat-firebase-emulators-9150.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9199,fork,reuseaddr TCP:firebase-emulators:9199 > /tmp/socat-firebase-emulators-9199.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:9299,fork,reuseaddr TCP:firebase-emulators:9299 > /tmp/socat-firebase-emulators-9299.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:8085,fork,reuseaddr TCP:firebase-emulators:8085 > /tmp/socat-firebase-emulators-8085.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:8090,fork,reuseaddr TCP:firebase-emulators:8090 > /tmp/socat-firebase-emulators-8090.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:5001,fork,reuseaddr TCP:firebase-emulators:5001 > /tmp/socat-firebase-emulators-5001.log 2>&1 &
bg_pids+=("$!")
socat TCP-LISTEN:4000,fork,reuseaddr TCP:firebase-emulators:4000 > /tmp/socat-firebase-emulators-4000.log 2>&1 &
bg_pids+=("$!")

# access to backend using localhost (position retrieval unauthorized using insecure http frontend with google-chrome)
socat TCP-LISTEN:8080,fork,reuseaddr TCP:application-backend:8080 > /tmp/socat-backend-8080.log 2>&1 &
bg_pids+=("$!")

# access to debug backend using localhost (position retrieval unauthorized using insecure http frontend with google-chrome)
socat TCP-LISTEN:8888,fork,reuseaddr TCP:devcontainer:8888 > /tmp/socat-devcontainer-8888.log 2>&1 &
bg_pids+=("$!")

# Developer-friendly aliases
alias gs='git status'
alias gb='git branch'
alias gd='git diff'
alias gl='git log --oneline --graph --decorate'
alias flutter-clean='flutter clean && rm -rf .dart_tool .packages pubspec.lock build'
alias run-android='flutter run -d android'

# Some linux distibution like Codespaces are requiring this additionnaly to the docker group addition.
sudo chown root:kvm /dev/kvm
sudo chmod 660 /dev/kvm

echo fs.inotify.max_user_watches=524288 |  sudo tee -a /etc/sysctl.conf; sudo sysctl -p

if [ "${LOCAL_ANDROID_STUDIO_APPIUM_EMULATOR_AVD_ON_START}" = "true" ]; then
    appium-emulator-avd.sh &
    bg_pids+=("$!")
fi

if [ "${LOCAL_ANDROID_STUDIO_ON_START}" = "true" ]; then
    android-studio.sh &
    bg_pids+=("$!")
fi

exec "$@"
