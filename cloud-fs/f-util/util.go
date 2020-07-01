package f_util

import "os"

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
