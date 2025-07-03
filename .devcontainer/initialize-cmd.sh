#!/bin/bash

# This script is run on the host as devcontainer 'initializeCommand' 
# (cf. https://containers.dev/implementors/json_reference/#lifecycle-scripts)

set -eu

trap "echo Exited with code $?." EXIT

# Create default local .env file with minimum required values to start.
localDotenvFile=${PWD}/.env

[ -f $localDotenvFile ] || cat <<'EOF' > $localDotenvFile
######################################################################## 
# After setting up values in this file, rebuild the local containers.  #
########################################################################
#  
#------------------------------------------------------- 
# Please set the values in this section according to your personnal settings.
# 
# Trigger the local project display name in Docker Compose.
COMPOSE_PROJECT_NAME=moov-dev-local
# 
# Make sure to set the correct values for using your personnal credentials IAM permissions. 
PROJECT_USER=user-to-setup
# 
# Can set 'MAKE_DEV_ON_START=false' to restart only the 'dev' container (skip 'make dev' in container 'dev' docker-compose command).
MAKE_DEV_ON_START=true
# 
# Android Studio (openbox - x11vnc - Xvfb)
LOCAL_ANDROID_STUDIO_ON_START=true
# 
# Set to match your screen resolution (e.g. if you are using the GUI from docker compose android-studio container).
# DISPLAY_RESOLUTION=680x1440
#
# Required if runnind E2E tests (application/tests)
LOCAL_ANDROID_STUDIO_APPIUM_EMULATOR_AVD_ON_START=true
#
# Wether to currently run the local application tests on start.
# If set to 'true', the local application tests will be run on start.
MAKE_LOCAL_APPLICATION_TESTS_RUN_ON_START=true
# 
#------------------------------------------------------- 
# The following variables are used with the local development environment.
# 
DEV_GOOGLE_CLOUD_PROJECT_ID=${GOOGLE_CLOUD_PROJECT_ID}
DEV_GOOGLE_IDP_OAUTH_KEY_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/google-idp-oauth-key/versions/latest
DEV_GOOGLE_IDP_OAUTH_CLIENT_ID_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/google-idp-oauth-client-id/versions/latest
DEV_STRIPE_KEY_SECRET_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/stripe-key/versions/latest
# 
BUILDER_IMAGE=europe-west1-docker.pkg.dev/${DEV_GOOGLE_CLOUD_PROJECT_ID}/docker-repository-public/${DEV_GOOGLE_CLOUD_PROJECT_ID}:builder-latest
FIREBASE_ADMINSDK_SERVICEACCOUNT_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/firebase-adminsdk-service-account-key/versions/latest
FIREBASE_PROJECT_ID=${DEV_GOOGLE_CLOUD_PROJECT_ID}

LOCAL_FIREBASE_EMULATORS_PUBSUB_VEGETABLE_IMAGES_VALIDATED_BACKEND_SUBSCRIPTION=vegetable-images-validated-backend
LOCAL_FIREBASE_EMULATORS_PUBSUB_VEGETABLE_IMAGES_CREATED_TOPIC=vegetable-images-created

UI_CONFIG_FIREBASE_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/firebase-config-web/versions/latest
UI_CONFIG_GOOGLEMAPS_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/${PROJECT_USER}-googlemaps-web-api-key/versions/latest

FIREBASE_STORAGE_PUBLIC_PREFIX=https://firebasestorage.googleapis.com/v0/b/${DEV_GOOGLE_CLOUD_PROJECT_ID}.appspot.com/o
CDN_PUBLIC_PREFIX=https://cdn.mon-backend.com  # ton CDN public GCS
# 
#--------------------------------------------------------
# ! Should not configure this section !^
#
# The following variables are used for propagating the containers
# configurations between them each others selves.
# 
ANDROID_HOST=android-studio
APPLICATION_BACKEND_URL=http://application-backend:8080
APPLICATION_BACKEND_DEBUG_URL=http://application-backend:8888
CLARINET_RPC=http://clarinet-devnet:20443
FIREBASE_AUTH_EMULATOR_HOST=firebase-emulators:9099
FIREBASE_DATABASE_EMULATOR_HOST=firebase-emulators:9000
FIREBASE_STORAGE_EMULATOR_HOST=firebase-emulators:9199
FIREBASE_PUBSUB_EMULATOR_HOST=firebase-emulators:8085
FIRESTORE_EMULATOR_HOST=firebase-emulators:8090
VAULT_ADDR=http://vault-dev:8200
VAULT_DEV_ROOT_TOKEN_ID=root
VAULT_DEV_LISTEN_ADDRESS=http://vault-dev:8200
STRIPE_KEY_PUBLISHABLE_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/stripe-key/versions/latest
STRIPE_KEY_SECRET_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/stripe-key/versions/latest
# 
# ! Should not configure this section !
#---------------------------------------------------------

EOF

# Vscode
workspaceFile=${PWD}/vscode.code-workspace
[ -f $workspaceFile ] || cat <<'EOF' > $workspaceFile
{
  "folders": [
    {
      "path": ".",
      "name": "Project"
    },
    {
      "name": "Devcontainer",
      "path": ".devcontainer"
    },
    {
      "name": "Local - Builder",
      "path": "local"
    },
    {
      "name": "Local - Firebase Emulators",
      "path": "firebase-emulators"
    },
    {
      "name": "Local - Android Studio",
      "path": "android-studio"
    },
    {
      "name": "Local - Vault",
      "path": "vault-dev"
    },
    {
      "name": "Local - Clarinet",
      "path": "clarinet-devnet"
    },
  ],
  "settings": {}
}
EOF

backendLaunchDebug=${PWD}/application/backend/.vscode/launch.json
if [ ! -f $backendLaunchDebug ] ;  then
mkdir -p $(dirname $backendLaunchDebug)
cat <<'EOF' > $backendLaunchDebug
{
    // Utilisez IntelliSense pour en savoir plus sur les attributs possibles.
    // Pointez pour afficher la description des attributs existants.
    // Pour plus d'informations, visitez : https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "env": {
                "PORT": "8888",
                "GOOGLE_APPLICATION_CREDENTIALS": "../../infra/environments/dev/gcloud-credentials.json",
                "UI_CONFIG_FIREBASE_SECRET_ID": "projects/${GOOGLE_CLOUD_PROJECT_ID}/secrets/firebase-config-web/versions/latest",
                "UI_CONFIG_GOOGLEMAPS_SECRET_ID": "projects/${GOOGLE_CLOUD_PROJECT_ID}/secrets/googlemaps-web-api-key/versions/latest",
                "STRIPE_KEY": "projects/${GOOGLE_CLOUD_PROJECT_ID}/secrets/stripe-key/versions/latest",
                "FIREBASE_PROJECT_ID": "${GOOGLE_CLOUD_PROJECT_ID}",
                "GCLOUD_PROJECT_ID": "${GOOGLE_CLOUD_PROJECT_ID}",
                "FRONTEND_BUILD_DIR": "../frontend/build",
                "FRONTEND_PUBLIC_DIR": "../frontend/public",
                "UI_JAVASCRIPT_SOURCE_FILE": "../frontend/build/bundle.js",
                "FIRESTORE_EMULATOR_HOST": "localhost:8090",
                "VAULT_ADDR": "http://localhost:8200",
                "VAULT_TOKEN": "root",
                "VAULT_MIN_USER_RECOVERY_KEY_ROTATION_INTERVAL": "1s",
                "FIREBASE_AUTH_EMULATOR_HOST": "localhost:9099",
                "VEGETABLE_CREATED_IMAGES_MODERATOR_PUBSUB_TOPIC": "vegetable-images-created",
                "VEGETABLE_VALIDATED_IMAGES_BACKEND_PUBSUB_SUBSCRIPTION": "vegetable-images-validated-backend",
                "VEGETABLE_VALIDATED_IMAGES_CDN_PREFIX_URL": "https://validated-images-cdn-prefix-url",
            },
            "envFile": "${workspaceFolder}/../../.env",
        }
    ]
}
EOF
fi

mobileLaunchDebug=${PWD}/application/mobile/.vscode/launch.json
if [ ! -f $mobileLaunchDebug ] ;  then
mkdir -p $(dirname $mobileLaunchDebug)
cat <<'EOF' > $mobileLaunchDebug
{
    // Utilisez IntelliSense pour en savoir plus sur les attributs possibles.
    // Pointez pour afficher la description des attributs existants.
    // Pour plus d'informations, visitez : https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "mobile",
            "request": "launch",
            "type": "dart"
        },
        {
            "name": "mobile (debug mode)",
            "request": "launch",
            "type": "dart",
            "flutterMode": "debug",
            "args": [
              "--dart-define=APPLICATION_BACKEND_URL=http://10.0.2.2:8888",
              "--dart-define=FIREBASE_STORAGE_PUBLIC_PREFIX=http://10.0.2.2:9199/v0/b/${GOOGLE_CLOUD_PROJECT_ID}.firebasestorage.app/o",            ]
        },
        {
            "name": "mobile (profile mode)",
            "request": "launch",
            "type": "dart",
            "flutterMode": "profile",
            "args": [
              "--dart-define=APPLICATION_BACKEND_URL=http://10.0.2.2:8080",
              "--dart-define=FIREBASE_STORAGE_PUBLIC_PREFIX=http://10.0.2.2:9199/v0/b/${GOOGLE_CLOUD_PROJECT_ID}.firebasestorage.app/o",            ]
        },
        {
            "name": "mobile (release mode)",
            "request": "launch",
            "type": "dart",
            "flutterMode": "release",
            "args": [
              "--dart-define=APPLICATION_BACKEND_URL=http://10.0.2.2:8080",
              "--dart-define=FIREBASE_STORAGE_PUBLIC_PREFIX=http://10.0.2.2:9199/v0/b/${GOOGLE_CLOUD_PROJECT_ID}.firebasestorage.app/o",            ]
        }
    ]
}
EOF
fi

CONTAINERS_CACHE_DIR=${PWD}/.containers
mkdir -p ${CONTAINERS_CACHE_DIR}

# Cache of container 'dev'
mkdir -p ${CONTAINERS_CACHE_DIR}/dev

# Copy config from host files.
if [ -d ~/.emacs.d ]; then
    rsync -av ~/.emacs.d ${CONTAINERS_CACHE_DIR}/emacs
fi