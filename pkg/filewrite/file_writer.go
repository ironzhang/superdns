package filewrite

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ironzhang/tlog"
)

// FileWriter file writer
type FileWriter struct {
	TemporaryDir string
}

func NewFileWriter(temporaryDir string) *FileWriter {
	return &FileWriter{TemporaryDir: temporaryDir}
}

func (p *FileWriter) WriteFile(path string, data []byte) (err error) {
	// 创建文件目录
	dir := filepath.Dir(path)
	if err = mkdir(dir); err != nil {
		tlog.Errorw("mkdir file dir", "dir", dir, "error", err)
		return err
	}

	// 创建临时目录
	if err = mkdir(p.TemporaryDir); err != nil {
		tlog.Errorw("mkdir temporary dir", "dir", p.TemporaryDir, "error", err)
		return err
	}

	// 写入临时文件
	file := p.temporaryFile(path)
	if err = ioutil.WriteFile(file, data, 0666); err != nil {
		tlog.Errorw("write file", "file", file, "error", err)
		return err
	}

	// 将临时文件重命名为正式文件
	if err = os.Rename(file, path); err != nil {
		os.Remove(file)
		tlog.Errorw("rename file", "old_path", file, "new_path", path, "error", err)
		return err
	}

	return nil
}

func (p *FileWriter) temporaryFile(path string) string {
	filename := filepath.Base(path)
	filename = fmt.Sprintf("%s.tmp.%s.%d", filename, time.Now().Format("2006-01-02T15:04:05.999999"), rand.Int())
	return filepath.Join(p.TemporaryDir, filename)
}
