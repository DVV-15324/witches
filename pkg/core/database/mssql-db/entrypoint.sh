#!/bin/bash

/opt/mssql/bin/sqlservr &

echo "Waiting for SQL Server..."

for i in {1..30}
do
    /opt/mssql-tools18/bin/sqlcmd \
        -S localhost \
        -U sa \
        -P "$SA_PASSWORD" \
        -C \
        -Q "SELECT 1" && break

    echo "SQL Server is starting..."
    sleep 2
done

echo "Creating database if not exists..."

/opt/mssql-tools18/bin/sqlcmd \
    -S localhost \
    -U sa \
    -P "$SA_PASSWORD" \
    -C \
    -Q "IF DB_ID('$MSSQL_DATABASE') IS NULL CREATE DATABASE [$MSSQL_DATABASE];"

echo "Running init scripts..."