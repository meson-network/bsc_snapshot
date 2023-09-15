rm -f -R ./build
mkdir build


echo "Compiling Mac amd64 version"
DIR="darwin-amd64"  &&  GOOS=darwin GOARCH=amd64 go build -a -o "./build/bsc_snapshot_darwin" 

echo "Compiling Windows amd64 version"
DIR="windows-amd64"  &&  GOOS=windows GOARCH=amd64 go build -a -o "./build/bsc_snapshot.exe" 

echo "Compiling Linux amd64 version"
DIR="linux-amd64"  &&  GOOS=linux GOARCH=amd64 go build -a -o "./build/bsc_snapshot_linux_amd64" 

echo "Compiling Linux ARM64 version"
DIR="linux-arm64"  &&  GOOS=linux GOARCH=arm64 go build -a -o "./build/bsc_snapshot_linux_arm64" 