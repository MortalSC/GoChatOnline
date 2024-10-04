# 编译 client
Write-Host "Building client executable..."
cd client
go build -o client.exe client.go
Write-Host "Client built successfully: client/client.exe"
cd ..

# 编译 server
Write-Host "Building server executable..."
cd server
go build -o server.exe main.go server.go user.go
Write-Host "Server built successfully: server/server.exe"
cd ..
