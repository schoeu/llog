package util

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var cfg *Config

type singleConfig struct {
	SysInfo       bool     `yaml:"sys_info"`
	LogDir        []string `yaml:"log_path"`
	Exclude       []string `yaml:"exclude_lines"`
	Include       []string `yaml:"include_lines"`
	ExcludeFiles  []string `yaml:"exclude_files"`
	MaxBytes      int      `yaml:"max_bytes"`
	TailFiles     bool     `yaml:"tail_files"`
	ScanFrequency int      `yaml:"scan_frequency"`
	CloseInactive int      `yaml:"close_inactive"`
	Multiline     struct {
		Pattern  string
		MaxLines int `yaml:"max_lines"`
	}
}

type outputConfig struct {
	ApiServer struct {
		Enable bool
		Url    string
	}
	Elasticsearch struct {
		Enable   bool
		Host     []string
		Index    string
		Username string
		Password string
	}
}

type Config struct {
	Name     string
	MaxProcs int `yaml:"max_procs"`
	Input    []singleConfig
	Output   outputConfig
}

func InitCfg(p string) error {
	p = GetAbsPath(GetCwd(), p)

	data, err := ioutil.ReadFile(p)
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	return err
}

func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	ErrHandler(errors.New("config not init"))
	return nil
}
