#!/bin/bash

go test ./...
make
python3 integration_tests.py