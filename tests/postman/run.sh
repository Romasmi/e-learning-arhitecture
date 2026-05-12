#!/bin/bash
newman run \
  tests/postman/elearning.postman_collection.json \
  --environment tests/postman/elearning.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export tests/postman/results.json \
  --bail
