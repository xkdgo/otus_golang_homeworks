#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER hwuser WITH password 'hwpasswd';
    CREATE DATABASE calendar;
    GRANT ALL PRIVILEGES ON DATABASE calendar TO hwuser;
EOSQL