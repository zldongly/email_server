package conf

type Conf struct {
	Server *Server `yaml:"server"`
	Data   *Data   `yaml:"data"`
}
