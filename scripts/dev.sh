#!/usr/bin/env bash

air -c ./.air.toml & \
bun x tailwindcss -i ./web/css/style.css -o ./web/static/main.css --watch