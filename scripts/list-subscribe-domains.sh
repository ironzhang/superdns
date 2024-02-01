#!/bin/bash

Addr="127.0.0.1:1789"

curl -X GET "http://$Addr/superdns/agent/v1/api/list/subscribe/domains"
