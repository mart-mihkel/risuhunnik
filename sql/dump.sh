#!/usr/bin/env bash

DB_FILE="build/risuhunnik.db"
DUMPS="build/dumps"
DUMP="$DUMPS/$(date +%Y%m%d%H%M%S).sql"

mkdir -p $DUMPS
sqlite3 $DB_FILE .dump > $DUMP
