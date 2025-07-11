#!/bin/bash

set -eu

local_container_cache=${LOCAL_ANDROID_STUDIO_CONTAINER_CACHE:-${LOCAL_DIR:-${PWD}}/.containers/android-studio}
mkdir -p $local_container_cache

# Android Studio config
ANDROID_STUDIO_CONFIG=${HOME}/.config/Google
mkdir -p ${local_container_cache}/Google ${HOME}/.config
ln -s ${local_container_cache}/Google $ANDROID_STUDIO_CONFIG

# Drop lock eventually remaining after previous container run exit
rm -f ${ANDROID_STUDIO_CONFIG}/AndroidStudio2024.1/.lock

# # Android Studio config
ANDROID_STUDIO_CACHE=${HOME}/.cache/Google
mkdir -p ${local_container_cache}/.cache/Google ${HOME}/.cache
ln -s ${local_container_cache}/.cache/Google $ANDROID_STUDIO_CACHE

# # Android Studio config (local)
ANDROID_STUDIO_LOCAL_SHARE_CONFIG=${HOME}/.share/Google
mkdir -p ${local_container_cache}/localconfig ${HOME}/.share
ln -s ${local_container_cache}/localconfig $ANDROID_STUDIO_LOCAL_SHARE_CONFIG

# Android Studio Gradle config
ANDROID_STUDIO_GRADLE_CONFIG=${HOME}/.gradle
mkdir -p ${local_container_cache}/gradle
ln -s ${local_container_cache}/gradle $ANDROID_STUDIO_GRADLE_CONFIG

# Android Studio Java cache
ANDROID_STUDIO_JAVA_CONFIG=${HOME}/.java
mkdir -p ${local_container_cache}/java
ln -s ${local_container_cache}/java $ANDROID_STUDIO_JAVA_CONFIG

# Android Studio Dart config
ANDROID_STUDIO_FLUTTER_CONFIG=${HOME}/.dart-tool
mkdir -p ${local_container_cache}/.dart-tool
ln -s ${local_container_cache}/.dart-tool $ANDROID_STUDIO_FLUTTER_CONFIG

# Android Studio Flutter config
ANDROID_STUDIO_FLUTTER_CONFIG=${HOME}/flutter
rsync -av $ANDROID_STUDIO_FLUTTER_CONFIG ${local_container_cache}
rm -rf $ANDROID_STUDIO_FLUTTER_CONFIG
ln -s ${local_container_cache}/flutter $ANDROID_STUDIO_FLUTTER_CONFIG

# Android Studio Pub config
ANDROID_STUDIO_PUB_CACHE_CONFIG=${HOME}/.pub-cache
rsync -av $ANDROID_STUDIO_PUB_CACHE_CONFIG ${local_container_cache}
rm -rf $ANDROID_STUDIO_PUB_CACHE_CONFIG
ln -s ${local_container_cache}/.pub-cache $ANDROID_STUDIO_PUB_CACHE_CONFIG

# Android Studio SDK
ANDROID_STUDIO_HOME=${HOME}/Android
rsync -av $ANDROID_STUDIO_HOME ${local_container_cache}
rm -rf $ANDROID_STUDIO_HOME
ln -s ${local_container_cache}/Android $ANDROID_STUDIO_HOME

# Android Studio Cache
ANDROID_STUDIO_CACHE=${HOME}/.android
rsync -av $ANDROID_STUDIO_CACHE ${local_container_cache}
rm -rf $ANDROID_STUDIO_CACHE
ln -s ${local_container_cache}/.android $ANDROID_STUDIO_CACHE

# Android Studio Dart server config
ANDROID_STUDIO_DART_SERVER=${HOME}/.dartServer
mkdir -p ${local_container_cache}/.dartServer
rm -rf $ANDROID_STUDIO_DART_SERVER
ln -s ${local_container_cache}/.dartServer $ANDROID_STUDIO_DART_SERVER

# vscode-server config
ANDROID_STUDIO_VSCODE_SERVER=${HOME}/.vscode-server
mkdir -p ${local_container_cache}/.vscode-server
rm -rf $ANDROID_STUDIO_VSCODE_SERVER
ln -s ${local_container_cache}/.vscode-server $ANDROID_STUDIO_VSCODE_SERVER

# Bash history
BASH_HISTORY_PATH=${HOME}/.bash_history
mkdir -p ${local_container_cache}
rm -f $BASH_HISTORY_PATH
ln -s ${local_container_cache}/.bash_history $BASH_HISTORY_PATH
# Git config (optional but useful)
GIT_CONFIG_GLOBAL=${HOME}/.gitconfig
mkdir -p ${local_container_cache}/git
if [ -f "$GIT_CONFIG_GLOBAL" ]; then
  rsync -av "$GIT_CONFIG_GLOBAL" ${local_container_cache}/git/
  rm -f "$GIT_CONFIG_GLOBAL"
fi
ln -sf ${local_container_cache}/git/.gitconfig $GIT_CONFIG_GLOBAL

# CAUSES UNFORTUNATE START WITH INCONSISTENT ALREADY RUNNING STATE
# # Persist VS Code settings (optional)
# VSCODE_SETTINGS=${HOME}/.config/Code
# mkdir -p $LOCAL_ANDROID_STUDIO/vscode
# if [ -d "$VSCODE_SETTINGS" ]; then
#   rsync -av "$VSCODE_SETTINGS" $LOCAL_ANDROID_STUDIO/vscode/
#   rm -rf "$VSCODE_SETTINGS"
# fi
# ln -sf $LOCAL_ANDROID_STUDIO/vscode/Code $VSCODE_SETTINGS
