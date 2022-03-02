package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ApiReadTimeOut  int      `yaml:"apiReadTimeOut"`
	ApiWriteTimeOut int      `yaml:"apiWriteTimeOut"`
	ApiPort         int      `yaml:"apiPort"`
	EtcdEndpoints   []string `yaml:"etcdEndpoints"`
	EtcdDialTimeOut int      `yaml:"etcdDialTimeOut"`
}

var (
	G_config *Config
)

func InitConfig(filePath string) error {

	var (
		content []byte
		err     error
		conf    *Config
	)

	//读取配置文件
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}
	//反序列化json到Config
	err = yaml.Unmarshal(content, &conf)
	G_config = conf
	return err
}
