#!/bin/bash

go test ./... -bench=. -count=5 -benchmem -run=^# | tee new.txt