#!/bin/sh

set -eu

trap "echo Exited with code $?." EXIT

mkdir -p ${VAULT_DATA}
mkdir -p ${VAULT_AUDIT}
mkdir -p ${VAULT_CONFIG}

tee ${VAULT_CONFIG}/vault.hcl <<EOF
ui            = true
cluster_addr  = "http://127.0.0.1:8201"
api_addr      = "https://127.0.0.1:8200"
disable_mlock = true

storage "raft" {
  path    = "${PWD}/.vault"
  node_id = "127.0.0.1"
}

listener "tcp" {
  address         = "0.0.0.0:8200"
  cluster_address = "0.0.0.0:8201"
  tls_disable     = 1
}
EOF


exec "$@"
