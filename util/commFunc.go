package util

import (
	"fmt"
	"os"
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
