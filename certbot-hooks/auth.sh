#!/usr/bin/env bash
set -euo pipefail

if [[ -n "${CERTBOT_TOKEN+x}" ]]; then
  # HTTP challenge
  echo "$CERTBOT_VALIDATION" >"$CERT_ACME_CHALLENGE_WEB_ROOT/$CERTBOT_TOKEN"
  echo "✅ Setup HTTP challenge for $CERTBOT_DOMAIN."
else
  # DNS challenge
  cd ./baidu-bce/
  ./baidu-bce record "_acme-challenge.$CERTBOT_DOMAIN" TXT "$CERTBOT_VALIDATION" --description "certbot at $(date --rfc-3339=s)"
  cd -
  echo "💤 Wait a second for DNS propagation of $CERTBOT_DOMAIN."
  sleep 1
  echo "✅ Setup DNS challenge for $CERTBOT_DOMAIN."
fi
