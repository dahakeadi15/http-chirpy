#!/bin/bash

source .env

cd sql/schema

if [[ "$1" == "up" ]]; then
    goose postgres $DB_URL up
elif [[ "$1" == "down" ]]; then
    goose postgres $DB_URL down
else
    echo "usage: scripts/runMigrations.sh [up|down]"
fi
