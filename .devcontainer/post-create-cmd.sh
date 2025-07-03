#!/bin/bash

set -eu

trap "echo Exited with code $?." EXIT

DEV_CONTAINER_CACHE=${PWD}/.containers/dev
mkdir -p $DEV_CONTAINER_CACHE

ln -sfn ${DEV_CONTAINER_CACHE}/bash_history ~/.bash_history

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
LOCAL_VSCODE_USER_GLOBAL_STORAGE=${DEV_CONTAINER_CACHE}/vscode/userData/globalStorage
mkdir -p ${LOCAL_VSCODE_USER_GLOBAL_STORAGE}
# persist locally (gitignored)
ln -sf ${DEV_CONTAINER_CACHE}/vscode/userData $VSCODE_REMOTE_USER_DATA
# versionned folder for gpt chat logging (folder ${DEV_CONTAINER_CACHE}/genieai.chatgpt-vscode)
ln -sf ${DEV_CONTAINER_CACHE}/genieai.chatgpt-vscode ${LOCAL_VSCODE_USER_GLOBAL_STORAGE}/
fi
