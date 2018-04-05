#!/bin/sh
sh install.sh
mkdir -p bin/
cd bin/
go build -o api-managemet ../server.go
go build -o rabbit-listener ../rabbit-listener.go
cd ..