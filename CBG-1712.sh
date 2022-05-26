#!/bin/bash

set -eux -o pipefail

env DOCKER_SCAN_SUGGEST=false docker-compose up --build --abort-on-container-exit
