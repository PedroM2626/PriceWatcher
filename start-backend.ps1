# Start the backend server
echo "Starting PriceWatcher backend..."
cd $PSScriptRoot
$env:GO111MODULE="on"
$env:GOFLAGS="-mod=vendor"

go run cmd/pricewatcher/main.go
