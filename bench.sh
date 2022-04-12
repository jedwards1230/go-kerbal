#!/bin/bash

rm test-debug.log
go test ./... -bench=. -benchtime=2x -count=5 -benchmem -run=^# | tee old.txt