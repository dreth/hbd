#!/bin/bash

# Check if the user has provided the DB_TYPE parameter
if [ -z "$1" ]; then
  echo 'Error: Please provide the DB_TYPE parameter. Either psql or sqlite3' >&2
  exit 1
fi

# Check the DB_TYPE parameter
DB_TYPE=$1
export DB_TYPE=$DB_TYPE

# if the DB_TYPE is sqlite, check if the file 'hbd.db' is present, otherwise touch it
if [ $DB_TYPE == "sqlite" ]; then
  if [ ! -f "hbd.db" ]; then
    touch hbd.db
  fi
fi

# check if we have the sqlboiler binary 
if ! [ -x "$(command -v sqlboiler)" ]; then
  echo 'Error: sqlboiler is not installed.' >&2
  exit 1
fi

# depending on the dbtype check if we have the sqlboiler-psql or the sqlboiler-sqlite3 binary
if [ $DB_TYPE == "psql" ]; then
  if ! [ -x "$(command -v sqlboiler-psql)" ]; then
    echo 'Error: sqlboiler-psql is not installed.' >&2
    exit 1
  fi
elif [ $DB_TYPE == "sqlite" ]; then
  if ! [ -x "$(command -v sqlboiler-sqlite3)" ]; then
    echo 'Error: sqlboiler-sqlite3 is not installed.' >&2
    exit 1
  fi
else
  echo 'Error: Invalid DB_TYPE. Please provide either psql or sqlite' >&2
  exit 1
fi

# apply up migrations
migrate -database $DATABASE_URL -path ./migrations/$DB_TYPE up

# If it's sqlite, generate the models for sqlite using sqlboiler
if [ $DB_TYPE == "sqlite" ]; then
  sqlboiler sqlite3 --config .sqlboiler-sqlite.toml --output ./models
  exit 0
fi
elif [ $DB_TYPE == "psql" ]; then
  sqlboiler psql --config .sqlboiler-postgres.toml --output ./models
  exit 0
fi
