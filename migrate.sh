#!/bin/bash

# A script to apply migrations to database. Provide a migration file name in arguments.

PGPASSWORD=postgres psql --no-password -h localhost -U postgres -d stackoverflow_recommender -a -f "${1}"
