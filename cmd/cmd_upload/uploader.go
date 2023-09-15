package cmd_upload

import (
	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc_snapshot/src/uploader"
)

func Uploader(clictx *cli.Context) error {

	originDir, thread, bucketName, additional_path,
		accountId, accessKeyId, accessKeySecret, retry_times := ReadParam(clictx)

	return uploader.Upload_r2(originDir, thread, bucketName, additional_path,
		accountId, accessKeyId, accessKeySecret, retry_times)
}
