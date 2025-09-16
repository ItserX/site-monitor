package config

type CheckerConfig struct {
	Checker struct {
		Interval int `yaml:"interval"`
		Timeout  int `yaml:"timeout"`
	} `yaml:"checker"`
}

type PostgresConfig struct {
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}

type CrudConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}
