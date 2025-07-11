#!/bin/bash

# This script is run on the host as devcontainer 'initializeCommand' 
# (cf. https://containers.dev/implementors/json_reference/#lifecycle-scripts)

set -eu

trap "echo Exited with code $?." EXIT
localFirebaseEmulatorsStorageBucket=${FIREBASE_EMULATORS_STORAGE_BUCKET:-${GOOGLE_CLOUD_PROJECT_ID}.appspot.com}
localFirebaseEmulatorsStorageRules=${LOCAL_FIREBASE_EMULATORS_STORAGE_RULES:-./storage.rules}
localFirebaseEmulatorsFirestoreRules=${LOCAL_FIREBASE_EMULATORS_FIRESTORE_RULES:-./firestore.rules}
localFirebaseEmulatorsAuthFunctionsDirectory=${LOCAL_FIREBASE_EMULATORS_AUTH_FUNCTIONS_DIR:-../../application/firebase/functions}

# Create default local .env file with minimum required values to start.
localFirebaseEmulatorsConfigJson=${LOCAL_FIREBASE_EMULATORS_CONFIG_JSON:-${LOCAL_DIR}/firebase.json}
[ -f $localFirebaseEmulatorsConfigJson ] || cat <<EOF > $localFirebaseEmulatorsConfigJson
{
  "storage": {
    "host": "0.0.0.0",
    "port": 9199,
    "rules": "${localFirebaseEmulatorsStorageRules}",
    "storageBucket": "${localFirebaseEmulatorsStorageBucket}"
  },
  "firestore": {
    "rules": "${localFirebaseEmulatorsFirestoreRules}"
  },
  "functions": {
    "source": "${localFirebaseEmulatorsAuthFunctionsDirectory}",
    "ignore": [".git", "firebase-debug.log", "firebase-debug.*.log", "*.local"]
  },

  "emulators": {
    "ui": {
      "enabled": true,
      "host": "0.0.0.0",
      "port": 4000
    },
    "functions": {
      "host": "0.0.0.0",
      "port": 5001
    },
    "pubsub": {
      "host": "0.0.0.0",
      "port": 8085
    },
    "firestore": {
      "host": "0.0.0.0",
      "port": 8090
    },
    "auth": {
      "host": "0.0.0.0",
      "port": 9099
    },
    "singleProjectMode": true,
    "database": {
      "host": "0.0.0.0",
      "port": 9000
    },

    "storage": {
      "host": "0.0.0.0",
      "port": 9199,
      "rules": "${localFirebaseEmulatorsStorageRules}",
      "storageBucket": "${localFirebaseEmulatorsStorageBucket}"
    }
  },
  "auth": {
    "providers": ["password", "google.com", "facebook.com", "anonymous"]
  }
}


EOF
