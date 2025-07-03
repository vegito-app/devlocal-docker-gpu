#!/bin/bash

set -eu

LOCAL_ANDROID_STUDIO=${PROJECT_DIR}/.containers/android-studio
mkdir -p $LOCAL_ANDROID_STUDIO

# Android Studio config
ANDROID_STUDIO_CONFIG=${HOME}/.config/Google
mkdir -p $LOCAL_ANDROID_STUDIO/Google ${HOME}/.config
ln -s $LOCAL_ANDROID_STUDIO/Google $ANDROID_STUDIO_CONFIG

# Drop lock eventually remaining after previous container run exit
rm -f ${ANDROID_STUDIO_CONFIG}/AndroidStudio2024.1/.lock

# # Android Studio config
ANDROID_STUDIO_CACHE=${HOME}/.cache/Google
mkdir -p $LOCAL_ANDROID_STUDIO/.cache/Google ${HOME}/.cache
ln -s $LOCAL_ANDROID_STUDIO/.cache/Google $ANDROID_STUDIO_CACHE

# # Android Studio config (local)
ANDROID_STUDIO_LOCAL_SHARE_CONFIG=${HOME}/.share/Google
mkdir -p $LOCAL_ANDROID_STUDIO/localconfig ${HOME}/.share
ln -s $LOCAL_ANDROID_STUDIO/localconfig $ANDROID_STUDIO_LOCAL_SHARE_CONFIG

# Android Studio Gradle config
ANDROID_STUDIO_GRADLE_CONFIG=${HOME}/.gradle
mkdir -p $LOCAL_ANDROID_STUDIO/gradle
ln -s $LOCAL_ANDROID_STUDIO/gradle $ANDROID_STUDIO_GRADLE_CONFIG

# Android Studio Java cache
ANDROID_STUDIO_JAVA_CONFIG=${HOME}/.java
mkdir -p $LOCAL_ANDROID_STUDIO/java
ln -s $LOCAL_ANDROID_STUDIO/java $ANDROID_STUDIO_JAVA_CONFIG

# Android Studio Dart config
ANDROID_STUDIO_FLUTTER_CONFIG=${HOME}/.dart-tool
mkdir -p $LOCAL_ANDROID_STUDIO/.dart-tool
ln -s $LOCAL_ANDROID_STUDIO/.dart-tool $ANDROID_STUDIO_FLUTTER_CONFIG

# Android Studio Flutter config
ANDROID_STUDIO_FLUTTER_CONFIG=${HOME}/flutter
rsync -av $ANDROID_STUDIO_FLUTTER_CONFIG $LOCAL_ANDROID_STUDIO
rm -rf $ANDROID_STUDIO_FLUTTER_CONFIG
ln -s $LOCAL_ANDROID_STUDIO/flutter $ANDROID_STUDIO_FLUTTER_CONFIG

# Android Studio Pub config
ANDROID_STUDIO_PUB_CACHE_CONFIG=${HOME}/.pub-cache
rsync -av $ANDROID_STUDIO_PUB_CACHE_CONFIG $LOCAL_ANDROID_STUDIO
rm -rf $ANDROID_STUDIO_PUB_CACHE_CONFIG
ln -s $LOCAL_ANDROID_STUDIO/.pub-cache $ANDROID_STUDIO_PUB_CACHE_CONFIG

# Android Studio SDK
ANDROID_STUDIO_HOME=${HOME}/Android
rsync -av $ANDROID_STUDIO_HOME $LOCAL_ANDROID_STUDIO
rm -rf $ANDROID_STUDIO_HOME
ln -s $LOCAL_ANDROID_STUDIO/Android $ANDROID_STUDIO_HOME

# Android Studio Cache
ANDROID_STUDIO_CACHE=${HOME}/.android
rsync -av $ANDROID_STUDIO_CACHE $LOCAL_ANDROID_STUDIO
rm -rf $ANDROID_STUDIO_CACHE
ln -s $LOCAL_ANDROID_STUDIO/.android $ANDROID_STUDIO_CACHE

# Android Studio Dart server config
ANDROID_STUDIO_DART_SERVER=${HOME}/.dartServer
mkdir -p $LOCAL_ANDROID_STUDIO/.dartServer
rm -rf $ANDROID_STUDIO_DART_SERVER
ln -s $LOCAL_ANDROID_STUDIO/.dartServer $ANDROID_STUDIO_DART_SERVER

# vscode-server config
ANDROID_STUDIO_VSCODE_SERVER=${HOME}/.vscode-server
mkdir -p $LOCAL_ANDROID_STUDIO/.vscode-server
rm -rf $ANDROID_STUDIO_VSCODE_SERVER
ln -s $LOCAL_ANDROID_STUDIO/.vscode-server $ANDROID_STUDIO_VSCODE_SERVER

# Bash history
BASH_HISTORY_PATH=${HOME}/.bash_history
mkdir -p $LOCAL_ANDROID_STUDIO
rm -f $BASH_HISTORY_PATH
ln -s $LOCAL_ANDROID_STUDIO/.bash_history $BASH_HISTORY_PATH
# Git config (optional but useful)
GIT_CONFIG_GLOBAL=${HOME}/.gitconfig
mkdir -p $LOCAL_ANDROID_STUDIO/git
if [ -f "$GIT_CONFIG_GLOBAL" ]; then
  rsync -av "$GIT_CONFIG_GLOBAL" $LOCAL_ANDROID_STUDIO/git/
  rm -f "$GIT_CONFIG_GLOBAL"
fi
ln -sf $LOCAL_ANDROID_STUDIO/git/.gitconfig $GIT_CONFIG_GLOBAL

# CAUSES UNFORTUNATE START WITH INCONSISTENT ALREADY RUNNING STATE
# # Persist VS Code settings (optional)
# VSCODE_SETTINGS=${HOME}/.config/Code
# mkdir -p $LOCAL_ANDROID_STUDIO/vscode
# if [ -d "$VSCODE_SETTINGS" ]; then
#   rsync -av "$VSCODE_SETTINGS" $LOCAL_ANDROID_STUDIO/vscode/
#   rm -rf "$VSCODE_SETTINGS"
# fi
# ln -sf $LOCAL_ANDROID_STUDIO/vscode/Code $VSCODE_SETTINGS
