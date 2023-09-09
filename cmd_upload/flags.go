package cmd_upload

import "github.com/urfave/cli/v2"

func GetFlags() (allflags []cli.Flag) {
	return []cli.Flag{
		&cli.StringFlag{Name: "dir", Required: true},
		&cli.StringFlag{Name: "bucket_name", Required: true},
		&cli.StringFlag{Name: "additional_path", Required: false},
		&cli.StringFlag{Name: "thread", Required: false},
		&cli.StringFlag{Name: "account_id", Required: true},
		&cli.StringFlag{Name: "access_key_id", Required: true},
		&cli.StringFlag{Name: "access_key_secret", Required: true},
	}
}


//originDir := clictx.String("dir")
	// thread := clictx.Int("thread")
	// bucketName := clictx.String("bucket_name")
	// accountId := clictx.String("account_id")
	// accessKeyId := clictx.String("access_key_id")
	// accessKeySecret := clictx.String("access_key_secret")