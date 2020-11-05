package config

import (
	"errors"
	"io/ioutil"

	"github.com/schoeu/llog/util"
	"gopkg.in/yaml.v2"
)

var cfg *Config

type SingleConfig struct {
	LogDir        []string `yaml:"log_path"`
	Type          string
	Exclude       []string `yaml:"exclude_lines"`
	Include       []string `yaml:"include_lines"`
	ExcludeFiles  []string `yaml:"exclude_files"`
	MaxBytes      int      `yaml:"max_bytes"`
	TailFiles     bool     `yaml:"tail_files"`
	ScanFrequency int      `yaml:"scan_frequency"`
	CloseInactive int      `yaml:"close_inactive"`
	Fields        string
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
	Name          string
	MaxProcs      int  `yaml:"max_procs"`
	SysInfo       bool `yaml:"sys_info"`
	SysInfoDuring int  `yaml:"sys_info_during"`
	Input         []SingleConfig
	Output        outputConfig
	SnapShot      struct {
		Enable         bool
		SnapshotDir    string `yaml:"snapshot_dir"`
		SnapShotDuring int    `yaml:"snapshot_during"`
	}
}

func InitCfg(p string) error {
	p = util.GetAbsPath(util.GetCwd(), p)

	data, err := ioutil.ReadFile(p)
	cfg = &Config{}
	err = yaml.Unmarshal(data, &cfg)
	return err
}

func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	util.ErrHandler(errors.New("config not init"))
	return nil
}
