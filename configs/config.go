package configs

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	DB     DBConfig
	Server ServerConf
}

var configPath = flag.String("conf", "config.yml", "config file path")

func Init() *Config {
	fmt.Println("init config")
	flag.Parse()

	conf := &Config{}
	log.Printf("load config file: %s", *configPath)
	// 读取配置文件
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("read %s fail, err: %v", *configPath, err)
	}

	// 解析yml文件
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatalf("parse config-local.yml fail, err: %v", err)
	}

	// 初始化fanbook配置
	//fanbook.InitConf(&conf.Fanbook)
	return conf
}
