#!/bin/bash

GENERATE_GROUPS="/Users/iron/workspace/src/golang/github.com/code-generator/generate-groups.sh"

$GENERATE_GROUPS all github.com/ironzhang/superdns/supercrd/clients github.com/ironzhang/superdns/supercrd/apis superdns.io:v1 --go-header-file ./boilerplate.go.txt

