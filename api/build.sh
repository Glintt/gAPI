#!/bin/sh
sh install.sh
mkdir -p bin/
cd bin/


platforms=("windows/amd64" "windows/386" "darwin/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name_server='api-'$GOOS'-'$GOARCH
    output_name_listener='api-listener-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name_server ../server.go
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name_listener ../rabbit-listener.go

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
cd ..

