#! /bin/sh

go test ./registry/database -bench=. -count 5 -run=^#
