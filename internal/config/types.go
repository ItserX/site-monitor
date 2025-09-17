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

	Prometheus struct {
		PushgatewayURL string `yaml:"pushgateway_url"`
	} `yaml:"prometheus"`
}

type CrudConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}

type AlertConfig struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`

	Telegram struct {
		BotToken string `yaml:"bot_token"`
		ChatID   string `yaml:"chat_id"`
	} `yaml:"telegram"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}
