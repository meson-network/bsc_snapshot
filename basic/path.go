package basic

import (
	"path/filepath"

	"github.com/coreservice-io/utils/path_util"
)

var WORK_DIR string

func AbsPath(rel_paths ...string) string {
	return filepath.Join(WORK_DIR, filepath.Join(rel_paths...))
}

// 1.if input is an abs path just return it if exist return true
// 2.if input is a relative path to working dir return abs_path if exist return true
func PathExist(abs_or_rel_path string) (string, bool, error) {

	//check abs path or relative path
	is_abs := filepath.IsAbs(abs_or_rel_path)
	if is_abs {
		//check if abs path exist
		abs_exist, err := path_util.AbsPathExist(abs_or_rel_path)
		if err != nil {
			return "", false, err
		}
		if abs_exist {
			return abs_or_rel_path, true, nil
		} else {
			return abs_or_rel_path, false, nil
		}
	} else {
		//rel to working directory
		abs_path := filepath.Join(WORK_DIR, abs_or_rel_path)
		exist, err := path_util.AbsPathExist(abs_path)
		if err != nil {
			return "", false, err
		}
		if exist {
			return abs_path, true, nil
		} else {
			return abs_path, false, nil
		}
	}
}
