#!/bin/bash

#seq 1 100 | parallel -j100
newman run \
  tests/postman/elearning.postman_collection.json \
  --environment tests/postman/elearning.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export tests/postman/results.json \
  --bail
