package cmd_download

import "github.com/urfave/cli/v2"

func GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "file_config", Required: true},
		&cli.StringFlag{Name: "thread", Required: false},
		&cli.IntFlag{Name: "retry_times", Required: false},
		&cli.BoolFlag{Name: "no_resume", Required: false},
	}
}

func ReadParam(clictx *cli.Context) (string, int, int, bool) {
	jsonConfigAddress := clictx.String("file_config")
	thread := clictx.Int("thread")
	retry_times := clictx.Int("retry_times")
	noResume := clictx.Bool("no_resume")

	return jsonConfigAddress, thread, retry_times, noResume
}
