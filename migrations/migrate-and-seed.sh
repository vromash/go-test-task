#!/bin/bash

goose postgres "user=postgres password=postgres host=postgres dbname=go-test-task-db sslmode=disable" up -dir ./migrations
goose postgres "user=postgres password=postgres host=postgres dbname=go-test-task-db sslmode=disable" up -dir ./migrations/seed --no-versioning