package main

import (
	. "es_test/conf"
	"es_test/wire"
	"fmt"

	"github.com/spf13/viper"
)

type UserSvr struct {
	conf Config
}

func (u *UserSvr) init() {
	viper.SetConfigName("config")
	//viper.SetConfigType("ini")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(err)
		}
		panic(err)
	}

	u.conf.Elastic.Host = viper.GetString("elastic.host")
	u.conf.Elastic.Port = viper.GetInt("elastic.port")
	u.conf.Elastic.Author = viper.GetString("elastic.author")
	u.conf.Elastic.Project = viper.GetString("elastic.project")
	fmt.Println(u.conf.Elastic)
}

func (u *UserSvr) Run() {
	handler := wire.InitializeHandler(&u.conf)
	handler.Run()
}
