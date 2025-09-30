#!/bin/bash

set -e
source .env

goose_cmd() {
  check_db $2

  if [ "$1" == "create" ] && [ -z "$3" ]; then
    echo "Please provide a migration name for create."
    exit 1
  fi

  if [ "$1" == "create" ]; then
    GOOSE_DRIVER="$GOOSE_DRIVER" GOOSE_DBSTRING="$DB_URL" \
      goose -dir internal/db/migrations "$1" "$3" sql
    exit 0
  fi

  GOOSE_DRIVER="$GOOSE_DRIVER" GOOSE_DBSTRING="$DB_URL" \
    goose -dir internal/db/migrations "$1"
}

check_db() {
  local env=$1
  if [ "$env" == "PROD" ]; then
      DB_URL=$PROD_DBSTRING
    elif [ "$env" == "DEV" ]; then
      DB_URL=$DEV_DBSTRING
    else
      echo "Please specify PROD or DEV as the second argument."
      exit 1
    fi
}

case "$1" in
  up|down|status)
   goose_cmd $1 $2
  ;;
  create)
    goose_cmd $1 $2 $3
  ;;
esac