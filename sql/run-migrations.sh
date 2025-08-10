#!/usr/bin/env bash

set -e

DB_FILE="build/risuhunnik.db"
LAST_FILE="build/last-migration"

touch $LAST_FILE

LAST_MIGRATION=$(cat $LAST_FILE)

if [[ -z $LAST_MIGRATION ]]; then
    LAST_DATE=0
else
    LAST_DATE=${LAST_MIGRATION:4:-4}
fi

echo "last migration at $LAST_DATE"

for MIGRATION in sql/*.sql; do
    MIGRATION_DATE=${MIGRATION:4:-4}
    if [[ MIGRATION_DATE -le LAST_DATE ]]; then
        echo "skipping $MIGRATION"
        continue
    fi

    echo "executing $MIGRATION"
    sqlite3 $DB_FILE < $MIGRATION
    echo $MIGRATION > $LAST_FILE
done
