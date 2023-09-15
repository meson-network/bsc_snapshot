rm -f -R ./build
mkdir build


echo "Compiling Mac amd64 version"
DIR="darwin-amd64"  &&  GOOS=darwin GOARCH=amd64 go build -a -o "./build/${DIR}/bsc_snapshot" 

echo "Compiling Windows amd64 version"
DIR="windows-amd64"  &&  GOOS=windows GOARCH=amd64 go build -a -o "./build/${DIR}/bsc_snapshot.exe" 

echo "Compiling Linux amd64 version"
DIR="linux-amd64"  &&  GOOS=linux GOARCH=amd64 go build -a -o "./build/${DIR}/bsc_snapshot" 

echo "Compiling Linux ARM64 version"
DIR="linux-arm64"  &&  GOOS=linux GOARCH=arm64 go build -a -o "./build/${DIR}/bsc_snapshot" 