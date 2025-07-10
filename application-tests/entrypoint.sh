#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

APPLICATION_TESTS_CONTAINER_CACHE=${HOST_PWD}/.containers/e2e-tests
mkdir -p $APPLICATION_TESTS_CONTAINER_CACHE

# Bash history
cat <<'EOF' >> ~/.bashrc
export HISTSIZE=50000
export HISTFILESIZE=100000
EOF

# Python/pip cache
PIP_CACHE_DIR=${HOME}/.cache/pip
[ -d $PIP_CACHE_DIR ] && mv $PIP_CACHE_DIR ${PIP_CACHE_DIR}_back || true
mkdir -p ${APPLICATION_TESTS_CONTAINER_CACHE}/pip ${PIP_CACHE_DIR}
ln -sf ${APPLICATION_TESTS_CONTAINER_CACHE}/pip $PIP_CACHE_DIR

kill_jobs() {
    echo "Killing background jobs"
    for pid in "$${bg_pids[@]}"; do
        kill "$$pid"
        wait "$$pid" 2>/dev/null
    done
}

trap kill_jobs EXIT

cat << 'EOF' >> ~/.bashrc
alias rf='robot --outputdir ${LOCAL_APPLICATION_TESTS_DIR} tests/robot'
alias h='htop'
alias i='sudo iftop'
alias ll='ls -lha'
alias l='ls -lh'
alias la='ls -Ah'
alias lla='ls -lhA'
EOF

mkdir -p ${LOCAL_APPLICATION_TESTS_DIR}

cd ${LOCAL_APPLICATION_TESTS_DIR} && python3 -m http.server 8088 &
bg_pids+=($!)

# Start Stripe CLI listener and forward webhooks to backend
stripe listen --forward-to ${APPLICATION_BACKEND_URL}/paiement/webhook &
bg_pids+=($!)
stripe listen --forward-to ${APPLICATION_BACKEND_DEBUG_URL}/paiement/webhook &
bg_pids+=($!)

exec "$@"

# Needed with github Codespaces which can change the workspace mount specified inside docker-compose.
current_workspace=$PWD
if [ "$current_workspace" != "$HOST_PWD" ] ; then
    sudo ln -s $current_workspace $HOST_PWD 2>&1 || true
    echo "Linked current workspace $current_workspace to $HOST_PWD"
fi
