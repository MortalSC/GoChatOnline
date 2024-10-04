@echo off
echo Building client executable...
cd client
go build -o client.exe client.go message.go
if %errorlevel% neq 0 (
    echo Failed to build client.
    exit /b %errorlevel%
)
echo Client built successfully: client/client.exe
cd ..

echo Building server executable...
cd server
go build -o server.exe main.go server.go user.go message.go
if %errorlevel% neq 0 (
    echo Failed to build server.
    exit /b %errorlevel%
)
echo Server built successfully: server/server.exe
cd ..

echo Build completed.
pause
