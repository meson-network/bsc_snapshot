package cmd_upload

import "github.com/urfave/cli/v2"

func GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "dir", Required: true},
		&cli.StringFlag{Name: "bucket_name", Required: true},
		&cli.StringFlag{Name: "additional_path", Required: false},
		&cli.IntFlag{Name: "thread", Required: false},
		&cli.StringFlag{Name: "account_id", Required: true},
		&cli.StringFlag{Name: "access_key_id", Required: true},
		&cli.StringFlag{Name: "access_key_secret", Required: true},
		&cli.IntFlag{Name: "retry_times", Required: false},
	}
}

func ReadParam(clictx *cli.Context) (string, int, string, string,
	string, string, string, int) {

	originDir := clictx.String("dir")
	thread := clictx.Int("thread")
	bucketName := clictx.String("bucket_name")
	additional_path := clictx.String("additional_path")
	accountId := clictx.String("account_id")
	accessKeyId := clictx.String("access_key_id")
	accessKeySecret := clictx.String("access_key_secret")
	retry_times := clictx.Int("retry_times")

	return originDir, thread, bucketName, additional_path,
		accountId, accessKeyId, accessKeySecret, retry_times
}
