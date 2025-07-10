#!/bin/bash

set -eu

local_container_cache=${LOCAL_FIREBASE_EMULATORS_CACHE:-${LOCAL_DIR:-${PWD}}/.firebase-emulators}
mkdir -p $local_container_cache

# Firebase Emulators config
FIREBASE_EMULATORS_CACHE=${HOME}/.cache
mkdir -p $local_container_cache/.cache
ln -s $local_container_cache/.cache $FIREBASE_EMULATORS_CACHE
