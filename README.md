# bsc-data-file-utils

This is a tool for splitting, uploading, and downloading large files, allowing users to easily split files and achieve multi-threaded uploads and downloads, significantly improving the speed of uploading and downloading large files.


## split file
Splitting the file will divide it into specified sizes and save it to the designated folder. Additionally, a 'files.json' file will be generated in the target folder to store information about the source file and the split files, making it convenient for various operations such as uploading and downloading.


split a large file and save to dest dir
```
 ./bsc-data-file-utils split \
 	--file=<file path> \
 	--dest=<to dir path> \
 	--size=<chunk size> \
    --thread=<thread quantity>
```
```
    --file   // <required> file path
    --size   // <required> each chunk size ex. 200m 
    --dest   // <optional> dest dir path ex. './dest'. default './dest'   
    --thread // <optional> thread quantity. default = cpu quantity
```


## set download endpoint
The 'endpoint' information in the 'config.json' file stores download endpoints, which allows automatic selection of download points when downloading files. Typically, multiple endpoints, in conjunction with multi-threaded downloads, can significantly increase download success rates and speed.
### add endpoints
add download endpoint
```
 ./bsc-data-file-utils endpoint add \
 	--config_path=<config file path> \
 	--endpoint=<endpoint url>
```
```
    --config_path   // <required> config file path
    --endpoint      // <required> url of endpoint to add, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

### remove endpoints
remove download endpoint
```
 ./bsc-data-file-utils endpoint remove \
 	--config_path=<config file path> \
 	--endpoint=<endpoint url>
```
```
    --config_path   // <required> config file path
    --endpoint      // <required> url of endpoint to remove, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

### set endpoints
set download endpoint, overwrite exist endpoints
```
 ./bsc-data-file-utils endpoint remove \
 	--config_path=<config file path> \
 	--endpoint=<endpoint url>
```
```
    --config_path   // <required> config file path
    --endpoint      // <required> url of endpoint to set, overwrite exist endpoints, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

### clear all endpoints
remove all endpoint
```
 ./bsc-data-file-utils endpoint remove \
 	--config_path=<config file path> \
```
```
    --config_path   // <required> config file path
```

### print exist endpoints
output exist endpoints
```
 ./bsc-data-file-utils endpoint remove \
 	--config_path=<config file path> \
```
```
    --config_path   // <required> config file path
```

## upload files

Upload the split files to storage. If the upload task is interrupted due to network or other reasons and needs to be resumed, typically, a comparison is made using MD5 for the files that have already been uploaded. Files with matching MD5 will not be re-uploaded.
### upload to cloudflare R2
upload to cloudflare R2 storage
```
 ./bsc-data-file-utils upload r2 \
 	--dir=<chunked file dir path> \
 	--bucket_name=<bucket name> \
 	--additional_path=<dir name> \
 	--account_id=<r2 account id> \
 	--access_key_id=<r2 access key id>  \
 	--access_key_secret=<r2 access key secret> \
    --thread=<thread quantity>  \
    --retry_times=<retry times>
```
```
    --dir               // <required> dir path to upload
 	--bucket_name       // <required> bucket name in r2
 	--additional_path   // <optional> dir name in bucket. default is "", means in root dir
 	--account_id        // <required> r2 account id
 	--access_key_id     // <required> r2 access key id
 	--access_key_secret // <required> r2 access key secret
    --thread            // <optional> thread quantity. default is 5
    --retry_times       // <optional> retry times limit when some file upload failed. default is 3
```

## download file
To download files, you need to provide 'files.json,' which is typically the file's download address (or it can also be a local file path). The download program will use the information in 'files.json' to perform multi-threaded downloads. During the download, the original source file is automatically reconstructed without the need for manual merging. Downloading supports resuming from breakpoints. If the download is interrupted due to network or other reasons, you simply need to rerun the download program to continue. After each small file is downloaded, an MD5 checksum is performed to ensure the integrity of the downloaded files.
```
 ./bsc-data-file-utils download \
 	--file_config=<json file url> \
    --thread=<thread quantity> \
    --retry_times=<retry times>
```
```
    --file_config   // <required> config file url
    --thread        // <optional> thread quantity. default is 5
    --retry_times   // <optional> retry times limit when some file download failed. default is 3
```


## file config struct
```
type FileConfig struct {
	RawFile         RawFileInfo       `json:"raw_file"`
	ChunkedFileList []ChunkedFileInfo `json:"chunked_file_list"`
	EndPoint        []string          `json:"end_point"`
}

type RawFileInfo struct {
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

type ChunkedFileInfo struct {
	FileName string `json:"file_name"`
	Md5      string `json:"md5"`
	Size     int64  `json:"size"`
	Offset   int64  `json:"offset"`
}
```