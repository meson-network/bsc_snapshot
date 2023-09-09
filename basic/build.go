package basic

import "errors"

const BUILD_MODE_RELEASE = "release"
const BUILD_MODE_DEBUG = "debug"

var mode = BUILD_MODE_RELEASE

func IsDebugMode() bool {
	return mode == BUILD_MODE_DEBUG
}

func IsReleaseMode() bool {
	return mode == BUILD_MODE_RELEASE
}

func SetMode(mode_ string) error {
	if mode_ != BUILD_MODE_RELEASE && mode_ != BUILD_MODE_DEBUG {
		return errors.New("mode not allowed :" + mode_)
	}
	mode = mode_
	return nil
}

func GetMode() string {
	return mode
}
