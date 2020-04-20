#!/usr/bin/env bash
set -eo pipefail

function finish {
	echo '[*] killing the api'
	kill "$(pgrep '[c]aliban')"
}
trap finish EXIT

echo '[*] starting the server'
./caliban &
sleep 2
echo '[*] testing the healthcheck'
curl -f "http://127.0.0.1:7080/health" 2> /dev/null | jq
echo '[*] GET /items'
curl -f "http://127.0.0.1:7080/items" 2> /dev/null | jq
echo '[*] POST /items/add'
curl -f -X POST -d '{"name":"item1", "description":"item 1"}' "http://127.0.0.1:7080/items/add" 2> /dev/null | jq
echo '[*] GET /items'
curl -f "http://127.0.0.1:7080/items" 2> /dev/null | jq
