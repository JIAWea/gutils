package gutils

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func CopyFile(from, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(to, input, 0777)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigDir() string {
	return filepath.Join(GetUserConfigDir(), "enputi")
}

func GetUserConfigDir() string {
	return filepath.Join(GetUserHomeDir(), ".config")
}

func GetUserHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func FileWrite(filename string, content string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = out.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func FileAppend(filename string, content string) error {
	out, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = out.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func FileRead(filename string) (content []byte, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return buf, err
}

func GetExecPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(execPath)
}

func GetPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return pwd
}

func GetFileModTime(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}
	}
	defer func() { _ = f.Close() }()

	fi, err := f.Stat()
	if err != nil {
		return time.Time{}
	}

	return fi.ModTime()
}

func FileReadByLine(fileName string, logic func(line string) error) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() { _ = fp.Close() }()

	r := bufio.NewReader(fp)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		err = logic(strings.TrimPrefix(string(line), "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
