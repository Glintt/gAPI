#!/bin/sh

cd $1

for d in */ ; do
    echo "Building plugins in: $d"
    cd $d
    for fi in *.go; do
	name=$(echo "$fi" | cut -f 1 -d '.')
	echo "Building $name"
       go build -buildmode=plugin -o $name.so $name.go
    done
    cd ..
done