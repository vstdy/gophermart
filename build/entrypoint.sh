#!/bin/sh
set -e

while ! nc -w 1 -zv postgres 5432; do sleep 1; done
./gophermart migrate

exec "$@"
