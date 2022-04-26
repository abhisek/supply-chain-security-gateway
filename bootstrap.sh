#!/bin/bash

set -ex

export ROOT_SUBJ="/C=IN/ST=KA/L=Bangalore/O=WeekendLabs/OU=DevOps/CN=weekend.labs/emailAddress=bofh@dev.null"

openssl req -nodes -new -x509 \
  -keyout pki/root.key -out pki/root.crt \
  -subj $ROOT_SUBJ

export PARTICIPATING_SERVICES="pdp tap dcs nats-server"

for svc in $PARTICIPATING_SERVICES; do
  mkdir -p pki/$svc

  openssl genrsa -out pki/$svc/server.key 2048

  # openssl req -x509 -key pki/$svc/server.key \
  #   -out pki/$svc/server.crt -sha256 -days 30 \
  #   -subj "/CN=$svc" \
  #   -addext "subjectAltName=DNS:$svc"

  openssl req -new -sha256 -key pki/$svc/server.key \
    -subj "/C=IN/ST=KA/O=WeekendLabs/CN=$svc" \
    -addext "subjectAltName = DNS:$svc" \
    -out pki/$svc/server.csr

  echo "[v3_req]\nsubjectAltName = DNS:$svc" > /tmp/$$-san.txt
  openssl x509 -req -in pki/$svc/server.csr \
    -CA pki/root.crt -CAkey pki/root.key -CAcreateserial \
    --extensions v3_req \
    -extfile /tmp/$$-san.txt \
    -out pki/$svc/server.crt -days 30 -sha256
  rm /tmp/$$-san.txt
done

# This is insecure but is needed for docker-compose
# In a production environment, we must use a cert manager instead of
# manually generating certificates

find ./pki -type f -exec chmod 644 {} \;

# Generate secrets
mysql_root_pass=$(openssl rand -hex 32)
cat > .env <<_EOF
MYSQL_ROOT_PASSWORD=$mysql_root_pass
MYSQL_DCS_DATABASE=vdb
MYSQL_DCS_USER=root
MYSQL_DCS_PASSWORD=$mysql_root_pass
_EOF
