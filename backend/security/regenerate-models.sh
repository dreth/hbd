#!/bin/sh
# Run from backend/. All database files below are disposable synthetic fixtures.
set -eu
hbd_generator_tmp=$(mktemp -d /tmp/hbd-model-generator.XXXXXXXX)
trap 'rm -rf "$hbd_generator_tmp"' EXIT HUP INT TERM
mkdir "$hbd_generator_tmp/bin"
export GOBIN="$hbd_generator_tmp/bin"
export PATH="$GOBIN:$PATH"
go install github.com/aarondl/sqlboiler/v4@v4.19.7
go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-sqlite3@v4.19.7
for migration in migrations/sqlite/*.up.sql; do sqlite3 "$hbd_generator_tmp/hbd.db" < "$migration"; done
export SQLITE3_DBNAME="$hbd_generator_tmp/hbd.db"
sqlboiler --replace 'main/07_relationship_to_one_eager.go.tpl;security/sqlboiler-templates/main/07_relationship_to_one_eager.go.tpl' --output models --pkgname models sqlite3
SQLITE3_DBNAME="$hbd_generator_tmp/hbd.db" go test ./...
