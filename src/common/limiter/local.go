package limiter

import (
	"time"

	"github.com/meson-network/bsc-data-file-utils/plugin/reference_plugin"
)

type limitInfo struct {
	Count_left        int
	Last_set_unixtime int64
}

// check key allowed to pass
func Allow(key string, duration_second int64, Count int) bool {

	if duration_second <= 0 {
		return true
	}

	lKey := "rateLimit:" + key
	value, _ := reference_plugin.GetInstance().Get(lKey)

	var limit_info *limitInfo
	nowTime := time.Now().UTC().Unix()

	allow := false

	if value == nil {

		allow = true

		limit_info = &limitInfo{
			Count_left:        Count,
			Last_set_unixtime: nowTime,
		}

		reference_plugin.GetInstance().Set(lKey, limit_info, duration_second*5+60)

	} else {
		limit_info = value.(*limitInfo)
		//if time past , add count
		if nowTime-limit_info.Last_set_unixtime >= duration_second {
			limit_info.Count_left = Count
			limit_info.Last_set_unixtime = nowTime
			allow = true
		} else {
			limit_info.Count_left--
			if limit_info.Count_left >= 0 {
				allow = true
			} else {
				allow = false
			}
		}
	}

	return allow
}
