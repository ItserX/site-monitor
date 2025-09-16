package config

type CheckerConfig struct {
	Checker struct {
		Interval int `yaml:"interval"`
		Timeout  int `yaml:"timeout"`
	} `yaml:"checker"`
}

type RedisConfig struct {
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}
