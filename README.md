> 参考链接：
> * https://github.com/asong2020/Golang_Dream/tree/master/code_demo
> * [go-ElasticSearch实战篇(二)](https://cloud.tencent.com/developer/article/2043026)
> * [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s?__biz=MzkzNjYxNTU3MQ==&mid=2247487180&idx=1&sn=ea9220b5f0a41540341eef3449c425fb&source=41&poc_token=HKviLmmjC3UBP0XceIyc1BYGQ_PFNdy3WzY2kRCz)

# 启动
```bash
go mod tidy
cd es_test
# 因为main.go 和 user_svr.go 的包都是 main
go run .
```