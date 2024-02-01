#!/bin/bash

Addr="127.0.0.1:1789"

ContentTypeHeader='Content-Type: application/json'

Data='{
	"Domains": ["example.app.com"],
	"TTL": "0s"
}'

curl -X POST "http://$Addr/superdns/agent/v1/api/subscribe/domains" --header "${ContentTypeHeader}" -d "$Data"
