#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

DEV_CONTAINER_CACHE=${PWD}/.containers/dev
mkdir -p $DEV_CONTAINER_CACHE

cat <<'EOF' >> ~/.bashrc
alias e="emacs"
EOF

# EMACS
EMACS_DIR=${HOME}/.emacs.d
[ -d $EMACS_DIR ] && mv $EMACS_DIR ${EMACS_DIR}_back
mkdir -p ${DEV_CONTAINER_CACHE}/emacs
ln -sf ${DEV_CONTAINER_CACHE}/emacs $EMACS_DIR

dev-entrypoint.sh "$@"