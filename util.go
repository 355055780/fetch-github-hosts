package main

import (
	"embed"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const (
	VERSION = 2.0
)

var (
	_debug    bool
	_execDir  string
	_logWrite io.Writer
)

func init() {
	_logWrite = os.Stdout
	_debug = os.Getenv("FETCH_GITHUB_HOST_DEBUG") != ""
	initAppExecDir()
}

func initAppExecDir() {
	if _debug {
		_execDir, _ = os.Getwd()
	} else {
		_exec, _ := os.Executable()
		_execDir = filepath.Dir(_exec)
	}
}

func IsDebug() bool {
	return _debug
}

func AppExecDir() string {
	return _execDir
}

func GetSystemHostsPath() string {
	switch runtime.GOOS {
	case Windows:
		return "C:/Windows/System32/drivers/etc/hosts"
	case Linux, Darwin:
		return "/etc/hosts"
	}
	return "/etc/hosts"
}

func PreCheckHasHostsRWPermission() (yes bool, err error) {
	_, err = syscall.Open(GetSystemHostsPath(), syscall.O_RDWR, 0655)
	if err != nil {
		if strings.Contains(err.Error(), "Access is denied") {
			err = nil
		}
		return
	}
	yes = true
	return
}

func ComposeError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return errors.New(msg + ": " + err.Error())
}

func GetExecOrEmbedFile(fs *embed.FS, filename string) (template []byte, err error) {
	exeDirFile := AppExecDir() + "/" + filename
	_, err = os.Stat(exeDirFile)
	if err == nil {
		template, err = ioutil.ReadFile(exeDirFile)
		return
	}
	template, err = fs.ReadFile(filename)
	return
}
