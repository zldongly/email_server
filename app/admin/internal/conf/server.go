package conf

type Server struct {
	Http *Http `yaml:"http"`
}

type Http struct {
	Addr string `yaml:"addr"`
}
