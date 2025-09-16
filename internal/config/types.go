package config

type CheckerConfig struct {
	Checker struct {
		Timeout  int    `yaml:"timeout"`
		Interval int    `yaml:"interval"`
		ApiURL   string `yaml:"api_url"`
	} `yaml:"checker"`

	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
	} `yaml:"kafka"`
}

type CrudConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}
