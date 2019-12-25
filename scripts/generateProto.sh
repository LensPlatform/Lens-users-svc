#!/usr/bin/env bash

for f in **/proto/*.proto; do \
    protoc --go_out=plugins=grpc:. $$f; \
    echo compiled: $$f; \
done