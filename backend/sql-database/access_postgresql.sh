#!/bin/sh
docker exec -it $1 psql -U $1
