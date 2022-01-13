#!/bin/bash

RETRIES=30
i=0
until docker-compose -f ./e2e/evm-optimism/docker-compose-nobuild.yml -f ./e2e/evm-optimism/docker-compose.e2e.yml logs relayer1 | grep -q "Starting relayer";
do
    sleep 3
    if [ $i -eq $RETRIES ]; then
        echo 'Timed out waiting for relayer'
        break
    fi
    echo 'Waiting for relayer...'
    ((i=i+1))
done
