#!/usr/bin/env bash
set -euo pipefail

if [[ -n "${CERTBOT_TOKEN+x}" ]]; then
  # HTTP challenge
  rm -f "$CERT_ACME_CHALLENGE_WEB_ROOT/$CERTBOT_TOKEN"
  echo "🧹 Clean up HTTP challenge for $CERTBOT_DOMAIN."
else
  # DNS challenge
  cd ./baidu-bce/
  ./baidu-bce forget "_acme-challenge.$CERTBOT_DOMAIN" TXT "$CERTBOT_VALIDATION"
  cd -
  echo "🧹 Clean up DNS challenge for $CERTBOT_DOMAIN."
fi
