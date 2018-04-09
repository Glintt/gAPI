#!/bin/sh
export API_MANAGEMENT_PORT=8080
export ELASTICSEARCH_HOST=localhost
export ELASTICSEARCH_PORT=9200
export SERVICEDISCOVERY_HOST=localhost
export SERVICEDISCOVERY_PORT=8080
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_USER=gapi
export RABBITMQ_PASSWORD=gapi
export RABBITMQ_QUEUE=gAPI-logqueue

sh install.sh
go run server.go