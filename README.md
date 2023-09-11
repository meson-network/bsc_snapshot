## bsc-data-file-utils


### split file
```
 ./bsc-data-file-utils split \
 	--file=<file path> \
 	--dest=<to dir path> \
 	--size=<chunk size> \
    --thread=<thread quantity>
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
    --thread=<thread quantity>
```

#### upload to aws s3



### download file
```
 ./bsc-data-file-utils download \
 	--file_config=<json file url>
    --thread=<thread quantity>
```