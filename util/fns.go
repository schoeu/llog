package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SysInfo      bool     `yaml:"sys_info"`
	LogDir       []string `yaml:"log_path"`
	Exclude      []string `yaml:"exclude_lines"`
	Include      []string `yaml:"include_lines"`
	ExcludeFiles []string `yaml:"exclude_files"`
	MaxBytes     int      `yaml:"max_bytes"`
	ApiServer    string   `yaml:"api_server"`
	TailFiles    bool     `yaml:"tail_files"`
	Multiline    struct {
		Pattern  string
		MaxLines int `yaml:"max_lines"`
	}
	Elasticsearch struct {
		Host     []string
		Protocal string
		Index    string
	}
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

func GetConfig(p string) (Config, error) {
	p = GetAbsPath(GetCwd(), p)

	c := Config{}
	data, err := ioutil.ReadFile(p)
	err = yaml.Unmarshal(data, &c)
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

func IsInclude(text string, regs []string) bool {
	for _, v := range regs {
		r := regexp.MustCompile(v)
		if r.MatchString(text) {
			return true
		}
	}
	return false
}
