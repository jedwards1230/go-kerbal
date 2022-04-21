#!/bin/bash

go test ./... -bench=. -benchtime=2x -count=5 -benchmem -run=^# | tee new.txt