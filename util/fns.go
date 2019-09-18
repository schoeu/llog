package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/satori/go.uuid"
)

type Config struct {
	LogDir    string
	NoSysInfo bool
	LogServer string
}

func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func GetHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return dir
}

func GetTempDir() string {
	return os.TempDir()
}

func ErrHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func GetAbsPath(base, p string) string {
	if base == "" {
		base = GetCwd()
	}
	if !path.IsAbs(p) {
		p = path.Join(base, p)
	}
	return p
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func UUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

func GetConfig(p string) (Config, error) {
	p = GetAbsPath(GetCwd(), p)

	c := Config{}
	data, err := ioutil.ReadFile(p)
	err = json.Unmarshal(data, &c)

	return c, err
}

func PathExist(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
