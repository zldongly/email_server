package conf

type Server struct {
	Kafka *Kafka `yaml:"kafka"`
}

type Kafka struct {
	Addrs []string `yaml:"addrs"`
}
