package conf

type Elastic struct {
	Host    string
	Port    int
	Author  string
	Project string
}

type Config struct {
	Elastic Elastic
}

var Conf Config
