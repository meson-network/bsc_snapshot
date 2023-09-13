rm -f -R ./build
mkdir build


# echo "Compiling Mac amd64 version"
# DIR="darwin-amd64"  &&  GOOS=darwin GOARCH=amd64 go build -a -o "./build/${DIR}/bsc-data-file-utils" 

# echo "Compiling Windows amd64 version"
# DIR="windows-amd64"  &&  GOOS=windows GOARCH=amd64 go build -a -o "./build/${DIR}/bsc-data-file-utils.exe" 

# echo "Compiling Windows arm64 version"
# DIR="windows-arm64"  &&  GOOS=windows GOARCH=arm64 go build -a -o "./build/${DIR}/bsc-data-file-utils.exe" 

echo "Compiling Linux amd64 version"
DIR="linux-amd64"  &&  GOOS=linux GOARCH=amd64 go build -a -o "./build/${DIR}/bsc-data-file-utils" 

# echo "Compiling Linux ARM32 version"
# DIR="linux-arm32"  &&  GOOS=linux GOARCH=arm go build -a -o "./build/${DIR}/bsc-data-file-utils" 

# echo "Compiling Linux ARM64 version"
# DIR="linux-arm64"  &&  GOOS=linux GOARCH=arm64 go build -a -o "./build/${DIR}/bsc-data-file-utils" 