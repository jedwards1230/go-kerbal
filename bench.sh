#!/bin/bash

rm test-debug.log
go test ./... -bench=. -benchtime=10x -benchmem -run=^#