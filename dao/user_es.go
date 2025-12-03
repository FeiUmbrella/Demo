package dao

import (
	"context"
	"es_test/conf"
	"es_test/model"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

// esRetryLimit 是重试次数
const esRetryLimit int = 3

// mappingTpl 表示es中的 mapping 结构
var mappingTpl string = `{
 "mappings":{
  "properties":{
   "id":     { "type": "long" },
   "username":   { "type": "keyword" },
   "nickname":   { "type": "text" },
   "phone":   { "type": "keyword" },
   "age":    { "type": "long" },
   "ancestral":  { "type": "text" },
   "identity":         { "type": "text" },
   "update_time":  { "type": "long" },
   "create_time":  { "type": "long" }
   }
  }
 }`

type UserES struct {
	client  *elastic.Client
	index   string
	mapping string
}

func NewUserES(client *elastic.Client, conf *conf.Config) (userEs *UserES) {
	index := fmt.Sprintf("%s_%s", conf.Elastic.Author, conf.Elastic.Project)
	userEs = &UserES{
		client:  client,
		index:   index,
		mapping: mappingTpl,
	}
	userEs.init()
	return
}

// init 初始化，主要是创建索引
func (userEs *UserES) init() {
	ctx := context.Background()

	exists, err := userEs.client.IndexExists(userEs.index).Do(ctx)
	if err != nil {
		fmt.Printf("userEs init exist failed err is %s\n", err.Error())
		return
	}

	// 创建的索引不存在，则创建对应索引
	if !exists {
		_, err = userEs.client.CreateIndex(userEs.index).Body(mappingTpl).Do(ctx)
		if err != nil {
			fmt.Printf("userEs init failed err is %s\n", err)
			return
		}
	}
}

// BatchAdd 批量添加
func (userEs *UserES) BatchAdd(ctx context.Context, user []*model.UserEs) error {
	var err error
	for _ = range esRetryLimit {
		if err = userEs.batchAdd(ctx, user); err != nil {
			fmt.Printf("es batchAdd failed err is %s\n", err)
			continue
		}
		return err
	}
	return err
}

// batchAdd 批量添加
func (userEs *UserES) batchAdd(ctx context.Context, user []*model.UserEs) error {
	req := userEs.client.Bulk().Index(userEs.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixMilli())
		u.CreateTime = uint64(time.Now().UnixMilli())

		// 创建单个文档的索引请求：
		// 1. 将用户ID转换为字符串作为文档ID
		// 2. 将用户对象u作为文档内容
		doc := elastic.NewBulkIndexRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		// 将单个文档请求添加到批量请求中
		req.Add(doc)
	}

	// req请求中没有要进行的操作
	if req.NumberOfActions() == 0 {
		return nil
	}
	// 执行批量请求到es中
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}

	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

// BatchUpdate 批量修改
func (userEs *UserES) BatchUpdate(ctx context.Context, user []*model.UserEs) error {
	var err error
	for _ = range esRetryLimit {
		if err = userEs.batchUpdate(ctx, user); err != nil {
			fmt.Printf("es batchUpdate failed err is %s\n", err)
			continue
		}
		return err
	}
	return err
}

// batchUpdate 批量修改
func (userEs *UserES) batchUpdate(ctx context.Context, user []*model.UserEs) error {
	req := userEs.client.Bulk().Index(userEs.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixMilli())

		// 更新单个文档的索引请求：
		// 1. 将用户ID转换为字符串作为文档ID
		// 2. 将用户对象u作为文档内容
		doc := elastic.NewBulkUpdateRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		// 将单个文档请求添加到批量请求中
		req.Add(doc)
	}

	// req请求中没有要进行的操作
	if req.NumberOfActions() == 0 {
		return nil
	}
	// 执行批量请求到es中
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}

	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

// BatchDel 批量删除
func (userEs *UserES) BatchDel(ctx context.Context, user []*model.UserEs) error {
	var err error
	for _ = range esRetryLimit {
		if err = userEs.batchDel(ctx, user); err != nil {
			fmt.Printf("es batchDel failed err is %s\n", err)
			continue
		}
		return err
	}
	return err
}

// batchDel 批量删除
func (userEs *UserES) batchDel(ctx context.Context, user []*model.UserEs) error {
	req := userEs.client.Bulk().Index(userEs.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixMilli())

		// 删除单个文档的索引请求：
		// 将用户ID转换为字符串作为文档ID进行查找删除
		doc := elastic.NewBulkDeleteRequest().Id(strconv.FormatUint(u.ID, 10))
		// 将单个文档请求添加到批量请求中
		req.Add(doc)
	}

	// req请求中没有要进行的操作
	if req.NumberOfActions() == 0 {
		return nil
	}
	// 执行批量请求到es中
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}

	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

func (userEs *UserES) Search(ctx context.Context, filter *model.EsSearch) ([]*model.UserEs, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(filter.MustQuery...)
	boolQuery.MustNot(filter.MustNotQuery...)
	boolQuery.Should(filter.ShouldQuery...)
	boolQuery.Filter(filter.Filters...)

	// 当should不为空，保证保证至少匹配should中的一项
	if len(filter.MustQuery) == 0 && len(filter.MustNotQuery) == 0 && len(filter.MustNotQuery) != 0 {
		boolQuery.MinimumShouldMatch("1")
	}

	service := userEs.client.Search().
		Index(userEs.index).
		Query(boolQuery).
		SortBy(filter.Sorters...).
		From(filter.From).
		Size(filter.Size)
	resp, err := service.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}

	userEss := make([]*model.UserEs, 0)
	for _, e := range resp.Each(reflect.TypeOf(&model.UserEs{})) {
		us := e.(*model.UserEs)
		userEss = append(userEss, us)
	}
	return userEss, nil
}
