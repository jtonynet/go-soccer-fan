#!/bin/sh

# Start the hot reload
CompileDaemon -build="go build -o /usr/src/app/bin/cli/main /usr/src/app/cmd/cli/main.go" -command="./bin/cli/main" &

# Keep the container running
tail -f /dev/null
