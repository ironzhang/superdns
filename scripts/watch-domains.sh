#!/bin/bash

Addr="127.0.0.1:1789"

ContentTypeHeader='Content-Type: application/json'

Data='{
	"Domains": ["example.app.com"]
}'

curl -X POST "http://$Addr/superdns/agent/v1/api/watch/domains" --header "${ContentTypeHeader}" -d "$Data"
