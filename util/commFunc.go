package util

import (
	"fmt"
	"os"
	"path"
)

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
		fmt.Println(err)
	}
}

func GetAbsPath(base, p string) string {
	if base == "" {
		base = GetCwd()
	}
	if !path.IsAbs(p) {
		p = path.Join(base, p)
	}
	fmt.Println(GetCwd())
	return p
}
