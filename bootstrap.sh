#!/bin/bash

set -ex

export ROOT_SUBJ="/C=IN/ST=KA/L=Bangalore/O=WeekendLabs/OU=DevOps/CN=weekend.labs/emailAddress=bofh@dev.null"

openssl req -nodes -new -x509 \
  -keyout pki/root.key -out pki/root.crt \
  -subj $ROOT_SUBJ

export SERVICES_IN_PKI="pdp tap nats"

for svc in $SERVICES_IN_PKI; do
  mkdir -p pki/$svc

  openssl genrsa -out pki/$svc/server.key 2048

  openssl req -new -sha256 -key pki/$svc/server.key \
    -subj "/C=IN/ST=KA/O=WeekendLabs/CN=$svc.weekend.labs" \
    -out pki/$svc/server.csr

  openssl x509 -req -in pki/$svc/server.csr \
    -CA pki/root.crt -CAkey pki/root.key -CAcreateserial \
    -out pki/$svc/server.crt -days 30 -sha256
done
