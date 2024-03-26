set dotenv-load

# List available recipes
@default:
    just --list

# Get a new certificate
get-cert DOMAINS *ARGS:
    certbot certonly --manual \
      --manual-auth-hook ./certbot-hooks/auth.sh \
      --manual-cleanup-hook ./certbot-hooks/cleanup.sh \
      --deploy-hook ./certbot-hooks/deploy.sh \
      -d '{{ DOMAINS }}' {{ ARGS }}

# Display information about certificates
show-cert *ARGS:
    certbot certificates {{ ARGS }}

# Call certbot (e.g. for debugging)
certbot *ARGS:
    certbot {{ ARGS }}
