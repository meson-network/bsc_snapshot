package cmd_split

import "github.com/urfave/cli/v2"

func GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "file", Required: true},
		&cli.StringFlag{Name: "dest", Required: false},
		&cli.StringFlag{Name: "size", Required: true},
		&cli.StringFlag{Name: "thread", Required: false},
	}
}

func ReadParam(clictx *cli.Context) (string, string, string, int) {
	originFilePath := clictx.String("file")
	destDir := clictx.String("dest")
	sizeStr := clictx.String("size")
	thread := clictx.Int("thread")

	return originFilePath, destDir, sizeStr, thread
}
