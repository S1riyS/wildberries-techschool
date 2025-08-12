#!/bin/bash
DSN=postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable

goose -dir "migrations" postgres "${DSN}" up -v