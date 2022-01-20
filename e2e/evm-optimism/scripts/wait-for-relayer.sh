#!/bin/bash

RETRIES=150
i=0
until docker-compose -f docker-compose-nobuild.yml -f docker-compose.e2e.yml logs relayer1 | grep -q "Starting relayer";
do
    sleep 5
    if [ $i -eq $RETRIES ]; then
        echo 'Timed out waiting for relayer'
        break
    fi
    echo 'Waiting for relayer...'
    ((i=i+1))
done
