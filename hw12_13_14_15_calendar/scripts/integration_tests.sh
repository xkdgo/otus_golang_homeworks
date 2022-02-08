#!/bin/bash

set -e

cd deployments
source tests.env

trap "docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose-tests.yaml down --rmi local &&\
 docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose-tests.yaml down -v" EXIT

docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose-tests.yaml up -d --scale integration-tests=0
docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose-tests.yaml run integration-tests
