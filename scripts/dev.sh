#!/usr/bin/env bash

# function cleanup() {
#   echo "Cleaning up..."
#   kill 0
# }
#
# trap cleanup SIGTERM SIGINT
# trap 'kill $(jobs -p)' EXIT
#
# bun x tailwindcss -i ./web/css/style.css -o ./web/static/main.css --watch &\
air -c ./.air.toml 

