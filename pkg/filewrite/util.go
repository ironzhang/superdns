package filewrite

import (
	"fmt"
	"os"
)

func mkdir(dir string) error {
	if fi, err := os.Stat(dir); err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("%q is not a dir", dir)
		}
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}
