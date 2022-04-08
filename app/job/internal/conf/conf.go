package conf

type Conf struct {
	Server *Server `yaml:"server"`
	Mail   *Mail   `yaml:"mail"`
	Data   *Data   `yaml:"data"`
}

type Mail struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
