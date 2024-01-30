#!/bin/sh

set -e

source /app/app.env

echo "start the app"
exec "$@"
