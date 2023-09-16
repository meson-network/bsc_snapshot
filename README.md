# BSC Snapshot

[![CodeFactor](https://www.codefactor.io/repository/github/meson-network/bsc_snapshot/badge)](https://www.codefactor.io/repository/github/meson-network/bsc_snapshot)
[![GitHub](https://img.shields.io/github/license/meson-network/bsc_snapshot)](https://github.com/meson-network/bsc_snapshot/blob/main/LICENSE)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/NetworkMeson)
<a href="https://discord.com/invite/z6YfSHDkmS">
    <img src="https://img.shields.io/badge/discord-join-7289DA.svg?logo=discord&longCache=true&style=flat" />
</a>
![GitHub forks](https://img.shields.io/github/forks/meson-network/bsc_snapshot)
![GitHub issues](https://img.shields.io/github/issues/meson-network/bsc_snapshot)
![GitHub all releases](https://img.shields.io/github/downloads/meson-network/bsc_snapshot/total)

This is a tool for splitting, uploading, and downloading large files, allowing users to easily split files and achieve multi-threaded uploads and downloads, significantly improving the speed of uploading and downloading large files.

Comparison with other download utils:

|  Util   | Speed  |
|  ----  | ----  |
| wget  | [#####--------------------------------------------------------]  35 MB/s |
| aria2c  | [#########################-----------------------------]  350 MB/s |
| bsc_snapshot  | [###############################################]  500 MB/s|


## How to use

1. Download this util and grant run permissions

* Linux 64bit
  
```text
wget -O bsc_snapshot "https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.1/bsc_snapshot_linux_amd64" && chmod +x ./bsc_snapshot
```

* Mac 64bit
  
```text
wget -O bsc_snapshot "https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.1/bsc_snapshot_darwin" && chmod +x ./bsc_snapshot
```

* Windows 64bit
  
```text
https://github.com/meson-network/bsc_snapshot/releases/download/v1.0.1/bsc_snapshot.exe
```

2. Start download

```text
 ./bsc_snapshot download --file_config=https://pub-51a999999b804ea79dd5bce1cb10c4e4.r2.dev/geth-20230907/files.json
```

param description:

```text
    --file_config   // <required> files.json url
    --thread        // <optional> thread quantity. default is 128
    --no_resume     // <optional> default is false, if set true, it will re-download file without resume
    --retry_times   // <optional> retry times limit when some file download failed. default is 5
```

To download files, you need to provide 'files.json,' which is typically the file's download address (or it can also be a local file path). The download program will use the information in 'files.json' to perform multi-threaded downloads. During the download, the original source file is automatically reconstructed without the need for manual merging. Downloading supports resuming from breakpoints. If the download is interrupted due to network or other reasons, you simply need to rerun the download program to continue. After each small file is downloaded, an MD5 checksum is performed to ensure the integrity of the downloaded files.

---

## For file deployment

### Step 1. split file

Splitting the file will divide it into specified sizes and save it to the designated folder. Additionally, a 'files.json' file will be generated in the target folder to store information about the source file and the split files, making it convenient for various operations such as uploading and downloading.

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

The 'endpoint' information in the 'files.json' stores download endpoints, which allows automatic selection of download points when downloading files. Typically, multiple endpoints, in conjunction with multi-threaded downloads, can significantly increase download success rates and speed.

The endpoint needs to be specified to a specific directory where files are stored, for example, if a file's download address is `https://yourdomain.com/bucket_dir/file.1`, then the endpoint should be set to `https://yourdomain.com/bucket_dir`.

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

Upload the split files to storage. If the upload task is interrupted due to network or other reasons and needs to be resumed, typically, a comparison is made using MD5 for the files that have already been uploaded. Files with matching MD5 will not be re-uploaded.

#### Upload to cloudflare R2

To upload files to Cloudflare R2, first, you need to create a bucket on R2 and obtain the 'account id', 'access key id', 'access key secret'.

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
