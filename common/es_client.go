package common

import (
	"fmt"
	"log"
	"os"

	"es_test/conf"

	"github.com/olivere/elastic/v7"
)

func NewEsClient(conf *conf.Config) *elastic.Client {
	url := fmt.Sprintf("http://%s:%d", conf.Elastic.Host, conf.Elastic.Port)
	client, err := elastic.NewClient(
		//elastic 服务地址
		elastic.SetURL(url),
		// 关闭嗅探器
		elastic.SetSniff(false),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		log.Fatal(err)
		//log.Fatalln("Failed to create elastic client")
	}
	return client
}
