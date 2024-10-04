#!/bin/bash

echo "Building client executable..."
cd client
GOOS=windows go build -o client client.go
if [ $? -ne 0 ]; then
    echo "Failed to build client."
    exit $?
fi
echo "Client built successfully: client/client"
cd ..

echo "Building server executable..."
cd server
GOOS=windows go build -o server main.go server.go user.go
if [ $? -ne 0 ]; then
    echo "Failed to build server."
    exit $?
fi
echo "Server built successfully: server/server"
cd ..

echo "Build completed."