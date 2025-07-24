#!/bin/bash

# Build the image for linux/amd64
docker build --platform linux/amd64 -t gimix/deltachat-bot .
