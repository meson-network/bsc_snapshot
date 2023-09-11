## bsc-data-file-utils


### split file
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


### upload files
#### upload to cloudflare R2
```
 ./bsc-data-file-utils upload \
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
 	--additional_path   // <optional> dir name in bucket. default is ""
 	--account_id        // <required> r2 account id
 	--access_key_id     // <required> r2 access key id
 	--access_key_secret // <required> r2 access key secret
    --thread            // <optional> thread quantity. default is 5
    --retry_times       // <optional> retry times limit when some file upload failed. default is 3
```



### download file
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