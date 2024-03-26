#!/usr/bin/env bash
set -euo pipefail

cd ./baidu-bce/
./baidu-bce upload "$CERTBOT_DOMAIN" "$RENEWED_LINEAGE"
cd -
echo "🚀 Upload the certificate to $CERTBOT_DOMAIN."
