#!/bin/bash

go test ./... -bench=. -benchtime=10x -benchmem -run=^#