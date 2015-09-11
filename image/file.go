package image

import (
	"os"
)

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
