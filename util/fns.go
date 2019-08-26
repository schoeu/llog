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
