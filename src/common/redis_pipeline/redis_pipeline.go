package redis_pipeline

import (
	"context"
	"time"

	"github.com/coreservice-io/job"
	"github.com/go-redis/redis/v8"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
)

const exec_count_limit = 100
const exec_interval_limit_millisec = 2500
const exec_thread_count = 4
const cmd_channel_limit = 20000

var last_exec_time_unixmilli = time.Now().UTC().UnixMilli()

type PipelineCmd struct {
	Ctx       context.Context
	Operation string
	Key       string
	Args      []interface{}
}

var cmdListChannel = make(chan *PipelineCmd, cmd_channel_limit)

func ScheduleRedisPipelineExec() {
	const jobName = "ScheduleRedisPipelineExec"

	for i := 0; i < exec_thread_count; i++ {
		job.Start(context.Background(), job.JobConfig{
			Name:          jobName,
			Job_type:      job.TYPE_PANIC_REDO,
			Interval_secs: 1,
			Process_fn: func(j *job.Job) {
				for {
					if len(cmdListChannel) < 100 && time.Now().UTC().UnixMilli()-last_exec_time_unixmilli < exec_interval_limit_millisec {
						time.Sleep(250 * time.Millisecond)
						continue
					}
					exec()
				}
			},
			On_panic: func(job *job.Job, panic_err interface{}) {
				basic.Logger.Errorln(panic_err)
			},
		}, nil)
	}
}

func exec() {

	last_exec_time_unixmilli = time.Now().UTC().UnixMilli()

	pl := redis_plugin.GetInstance().Pipeline()

outLoop:
	for i := 0; i < exec_count_limit; i++ {
		select {
		case cmd := <-cmdListChannel:
			switch cmd.Operation {
			case operation_Set:
				pl.Set(cmd.Ctx, cmd.Key, cmd.Args[0], cmd.Args[1].(time.Duration))

			case operation_ZAdd:
				z := []*redis.Z{}
				for _, v := range cmd.Args {
					z = append(z, v.(*redis.Z))
				}
				pl.ZAdd(cmd.Ctx, cmd.Key, z...)

			case operation_ZAddNX:
				z := []*redis.Z{}
				for _, v := range cmd.Args {
					z = append(z, v.(*redis.Z))
				}
				pl.ZAddNX(cmd.Ctx, cmd.Key, z...)

			case operation_HSet:
				pl.HSet(cmd.Ctx, cmd.Key, cmd.Args...)

			case operation_Expire:
				pl.Expire(cmd.Ctx, cmd.Key, cmd.Args[0].(time.Duration))

			case operation_ZRemRangeByScore:
				pl.ZRemRangeByScore(cmd.Ctx, cmd.Key, cmd.Args[0].(string), cmd.Args[1].(string))

			default:
				basic.Logger.Errorln("unsupported cmd:", cmd.Operation)
			}

		default:
			break outLoop
		}
	}

	if pl.Len() == 0 {
		return
	}

	_, err := pl.Exec(context.Background())
	if err != nil {
		basic.Logger.Errorln("exec pipeline error:", err)
		time.Sleep(5 * time.Second) //sleep a while for exe err
		return
	}
}
