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
LOCAL_ANDROID_STUDIO_APK_PATH=application/mobile/build/app/outputs/flutter-apk/app-release.apk
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

# Set this value tu 'host' to use accelerated GPU rendering in Android Studio.
# Set to 'swiftshader_indirect' to use software rendering if you are not using a GPU.
LOCAL_ANDROID_STUDIO_ANDROID_GPU_MODE=swiftshader_indirect

# Set this value to 'Pixel_8_Intel' or 'Pixel_6_Playstore' to use the corresponding AVD.
LOCAL_ANDROID_STUDIO_ANDROID_AVD_NAME=Pixel_6_Playstore

UI_CONFIG_FIREBASE_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/firebase-config-web/versions/latest
UI_CONFIG_GOOGLEMAPS_SECRET_ID=projects/${DEV_GOOGLE_CLOUD_PROJECT_ID}/secrets/${PROJECT_USER}-googlemaps-web-api-key/versions/latest

FIREBASE_STORAGE_PUBLIC_PREFIX=https://firebasestorage.googleapis.com/v0/b/${DEV_GOOGLE_CLOUD_PROJECT_ID}.appspot.com/o
CDN_PUBLIC_PREFIX=https://cdn.mon-backend.com  # ton CDN public GCS
# 
#--------------------------------------------------------
# ! Should not configure this section !
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
