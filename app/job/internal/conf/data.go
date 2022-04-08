package conf

type Data struct {
	Mongo *Mongo `yaml:"mongo"`
}

type Mongo struct {
	Uri string `yaml:"uri"`
	// 认证
	AuthSource string `yaml:"auth_source"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}
