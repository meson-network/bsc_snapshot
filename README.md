# BSC Snapshot

[![CodeFactor](https://www.codefactor.io/repository/github/meson-network/bsc_snapshot/badge)](https://www.codefactor.io/repository/github/meson-network/bsc_snapshot)
[![GitHub](https://img.shields.io/github/license/meson-network/bsc_snapshot)](https://github.com/meson-network/bsc_snapshot/blob/main/LICENSE)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/NetworkMeson)

This repo solves the problem of binance BSC chain sync.

Tools are provided for both user and maintainer to reduce the sync time approximately from ~20 hours to ~1 hour.

The tools can be configured easily to split upload and download a large snapshot file with multi service providers.

Multi threads downloading with multi endpoints significantly reduce the time cost of sync .

The service providers currently include meson.network which is a globally distributed files cache layer.

GreenField can be configured to work as the storage layer.

Comparison with other download utils:

| Util | Speed |
| ---- | ---- |
| wget | [#####--------------------------------------------------------] 35 MB/s |
| aria2c | [#########################-----------------------------] 350 MB/s |
| bsc_snapshot | [###############################################] 500 MB/s|


## [For user] How to use

  

1. Download this util and grant execution permissions

  

* Linux 64bit

```text

wget -O bsc_snapshot "https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.3/bsc_snapshot_linux_amd64" && chmod +x ./bsc_snapshot

```

  

* Mac 64bit

```text

wget -O bsc_snapshot "https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.3/bsc_snapshot_darwin" && chmod +x ./bsc_snapshot

```

  

* Windows 64bit

```text

https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.3/bsc_snapshot.exe

```

  

2. Start download

  

```text

./bsc_snapshot download --file_config=https://pub-51a999999b804ea79dd5bce1cb10c4e4.r2.dev/geth-20230907/files.json

```

  

param description:

  

```text
    --file_config // <required> files.json url
    --thread // <optional> thread quantity. default is 128
    --no_resume // <optional> default is false, if set true, it will re-download file without resume
    --retry_times // <optional> retry times limit when some file download failed. default is 5
```

  
The  --file_config 'files.json' is required which is a config file used to perform multi-threaded downloading. The original source file is automatically reconstructed without manual merging. Resuming from breakpoints and MD5 checksum are already included which ensures the efficiency integrity and safety.

  
  

---------------------------------------------------------------------------------
---------------------------------------------------------------------------------
---------------------------------------------------------------------------------
---------------------------------------------------------------------------------

  

## [For maintainer] 

  

### Step 1. split file

  

Splitting file tool helps to divide the source file into specified sizes and save it to the designated folder. Additionally, a 'files.json' file will be generated in the target folder which is a config file making it convenient for various operations such as uploading and downloading.


#### Split a large file to dest dir


```text
 ./bsc_snapshot split \
    --file=<file path> \
    --dest=<to dir path> \
    --size=<chunk size> \
    --thread=<thread quantity>
```

  

Param description:

  

```text
    --file   // <required> file path
    --size   // <required> each chunk size ex. 200m 
    --dest   // <optional> dest dir path ex. './dest'. default './dest'   
    --thread // <optional> thread quantity. default = cpu quantity
```

  

#### files.json Struct

  
```golang
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

  

### Step 2. set download endpoint

  
The 'endpoint' in the 'files.json' stores download source entries which can be selected during downloading process.
More endpoints redundancy working together with multi-threaded downloading will provide benefits of efficiency and safety.
The endpoint should be specified to a specific directory where files are stored, for example, if a file's download address is `https://yourdomain.com/bucket_dir/file.1`, then the endpoint should be set to `https://yourdomain.com/bucket_dir`.

  

#### Add endpoints

  

add download endpoint

  

```text
 ./bsc_snapshot endpoint add \
    --config_path=<files.json path> \
    --endpoint=<endpoint url>
```

  

param description:

  

```text
    --config_path   // <required> files.json path
    --endpoint      // <required> endpoint url to add, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

  

#### Remove endpoints

  

remove download endpoint

  

```text
 ./bsc_snapshot endpoint remove \
    --config_path=<files.json path> \
    --endpoint=<endpoint url>
```

  

param description:

  

```text
    --config_path   // <required> files.json path
    --endpoint      // <required> url of endpoint to remove, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

  

#### Set endpoints

  

set download endpoint, overwrite exist endpoints

  

```text
 ./bsc_snapshot endpoint set \
    --config_path=<files.json path> \
    --endpoint=<endpoint url>
```

  

param description:

  

```text
    --config_path   // <required> files.json path
    --endpoint      // <required> url of endpoint to set, overwrite exist endpoints, support multiple endpoint, ex. --endpoint=<url1> --endpoint=<url2>
```

  

#### Clear all endpoints

  

remove all endpoint

  

```text
 ./bsc_snapshot endpoint clear \
    --config_path=<files.json path> \
```

  

param description:

  

```text
    --config_path   // <required> files.json path
```

  

#### Print exist endpoints

  

output exist endpoints

  

```text
 ./bsc_snapshot endpoint print \
    --config_path=<files.json path> \
```

  

param description:

  

```text
    --config_path   // <required> files.json path
```

  

### Step 3. upload files to storage

Files will be checked using MD5 to prevent re-uploading if the upload task is interrupted due to network or other reasons.


#### Upload to cloudflare R2 


Before uploading files to Cloudflare R2 you need to create a bucket on R2 and obtain the 'account id', 'access key id', 'access key secret'.
  

```text
 ./bsc_snapshot upload r2 \
    --dir=<chunked file dir path> \
    --bucket_name=<bucket name> \
    --additional_path=<dir name> \
    --account_id=<r2 account id> \
    --access_key_id=<r2 access key id>  \
    --access_key_secret=<r2 access key secret> \
    --thread=<thread quantity>  \
    --retry_times=<retry times>
```

  

param description:

  

```text
    --dir               // <required> dir path to upload
    --bucket_name       // <required> bucket name in r2
    --additional_path   // <optional> dir name in bucket. default is "", means in bucket root dir
    --account_id        // <required> r2 account id
    --access_key_id     // <required> r2 access key id
    --access_key_secret // <required> r2 access key secret
    --thread            // <optional> thread quantity. default is 5
    --retry_times       // <optional> retry times limit when some file upload failed. default is 5
```
