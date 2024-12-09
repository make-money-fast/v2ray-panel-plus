package helpers

import "os"

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
