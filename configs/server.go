package configs

type ServerConf struct {
	Port            int  `yaml:"port"`
	ShutdownTimeout int  `yaml:"shutdownTimeout"` // 优雅停止服务的超时时间（秒）
	OpenCors        bool `yaml:"openCors"`        // 是否开启cors处理

}
