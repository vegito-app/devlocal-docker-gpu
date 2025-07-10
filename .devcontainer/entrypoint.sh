#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

# Local Container Cache
local_container_cache=${LOCAL_local_container_cache:-${LOCAL_DIR:-${PWD}}/.containers/dev}
mkdir -p $local_container_cache

cat <<'EOF' >> ~/.bashrc
alias e="emacs"
EOF

# EMACS local configuration persistence
# This allows you to persist your emacs configuration across container rebuilds.
EMACS_DIR=${HOME}/.emacs.d
[ -d $EMACS_DIR ] && mv $EMACS_DIR ${EMACS_DIR}_back
mkdir -p ${local_container_cache}/emacs
ln -sf ${local_container_cache}/emacs $EMACS_DIR

# BASH history
ln -sfn ${local_container_cache}/bash_history ~/.bash_history

# Vscode server/remote
VSCODE_REMOTE=${HOME}/.vscode-server

# Github Codespaces
if [ -v  CODESPACES ] ; then
    VSCODE_REMOTE=${HOME}/.vscode-remote
fi

# VSCODE User data
VSCODE_REMOTE_USER_DATA=${VSCODE_REMOTE}/data/User
if [ -d $VSCODE_REMOTE_USER_DATA ] ; then 
mv $VSCODE_REMOTE_USER_DATA ${VSCODE_REMOTE_USER_DATA}_back
LOCAL_VSCODE_USER_GLOBAL_STORAGE=${local_container_cache}/vscode/userData/globalStorage
mkdir -p ${LOCAL_VSCODE_USER_GLOBAL_STORAGE}
# persist locally (gitignored)
ln -sf ${local_container_cache}/vscode/userData $VSCODE_REMOTE_USER_DATA
# versionned folder for gpt chat logging (folder ${local_container_cache}/genieai.chatgpt-vscode)
ln -sf ${local_container_cache}/genieai.chatgpt-vscode ${LOCAL_VSCODE_USER_GLOBAL_STORAGE}/
fi

dev-entrypoint.sh "$@"