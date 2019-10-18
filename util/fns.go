package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/satori/go.uuid"
)

func GetCwd() string {
	dir, err := os.Getwd()
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
	if !filepath.IsAbs(p) {
		p = filepath.Join(base, p)
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

func IsInclude(text []byte, regs []string) bool {
	for _, v := range regs {
		r := regexp.MustCompile(v)
		if r.Match(text) {
			return true
		}
	}
	return false
}

func Recover() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}
