package filewrite

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestTemporaryFile(t *testing.T) {
	tests := []struct {
		dir  string
		path string
		file string
	}{
		{
			dir:  "./tmp",
			path: "./a",
			file: "tmp/a.tmp.",
		},
		{
			dir:  "/home/xiaoju/.service/disf/tmp",
			path: "/home/xiaoju/.service/disf/disf!biz-gs-tripcloud_agent/__config.json",
			file: "/home/xiaoju/.service/disf/tmp/__config.json.tmp.",
		},
	}
	for i, tt := range tests {
		w := FileWriter{TemporaryDir: tt.dir}
		file := w.temporaryFile(tt.path)
		if got, want := file, tt.file; !strings.HasPrefix(got, want) {
			t.Errorf("%d: temporaryFile: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: temporaryFile: got %v", i, got)
		}
	}
}

func TestWriteFile(t *testing.T) {
	content := []byte("hello, world\n")
	w := FileWriter{TemporaryDir: "./testdata/test_write_file/tmp"}
	paths := []string{
		"./testdata/test_write_file/a.json",
		"./testdata/test_write_file/a.json",
		"./testdata/test_write_file/a/a.json",
		"./testdata/test_write_file/a/b.json",
	}
	for i, path := range paths {
		if err := w.WriteFile(path, content); err != nil {
			t.Errorf("%d: WriteFile(%q): %v", i, path, err)
			continue
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("%d: ReadFile(%q): %v", i, path, err)
			continue
		}
		if got, want := data, content; !bytes.Equal(got, want) {
			t.Errorf("%d: content: got %q, want %q", i, got, want)
		}
	}
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
