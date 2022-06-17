package main

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"syscall"
)

const (
	VERSION = 1.0
)

var (
	_debug   bool
	_execDir string
)

func init() {
	_debug = os.Getenv("FETCH_GITHUB_HOST_DEBUG") != ""
	initAppExecDir()
}

func initAppExecDir() {
	if _debug {
		_execDir, _ = os.Getwd()
	} else {
		_execDir, _ = os.Executable()
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
		if runtime.GOOS == Windows {
			if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
				err = nil
			}
		}
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
