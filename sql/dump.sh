#!/usr/bin/env bash

DB_FILE="build/risuhunnik.db"
BACKUPS="dumps"
DUMP="$BACKUPS/$(date +%Y%m%d%H%M%S).sql" 

mkdir -p $BACKUPS
sqlite3 $DB_FILE .dump > $DUMP
