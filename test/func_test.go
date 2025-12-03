package test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"es_test/conf"
	"es_test/dao"
	"es_test/model"

	"github.com/olivere/elastic/v7"
)

// TestBatchUpdate 批量更新的单测
/*
使用真实配置构造 ES 客户端、读取并备份 ID 6/8 的原始文档后执行 dao.UserES.BatchUpdate，
最后断言昵称、手机号、身份字段已更新，并在 defer 中回滚，确保测试可重复
*/
func TestBatchUpdate(t *testing.T) {
	ctx := context.Background()
	cfg := &conf.Config{
		Elastic: conf.Elastic{
			Host:    "10.224.201.42",
			Port:    9200,
			Author:  "sf",
			Project: "es_test",
		},
	}

	client := newESTestClient(t, cfg)
	userES := dao.NewUserES(client, cfg)
	index := fmt.Sprintf("%s_%s", cfg.Elastic.Author, cfg.Elastic.Project)
	ids := []uint64{6, 8}
	original := make(map[uint64]*model.UserEs, len(ids))
	for _, id := range ids {
		original[id] = fetchUser(ctx, t, client, index, id)
	}

	defer func() {
		recoverUsers := []*model.UserEs{cloneUser(original[6]), cloneUser(original[8])}
		if err := userES.BatchUpdate(ctx, recoverUsers); err != nil {
			t.Fatalf("恢复原始数据失败: %v", err)
		}
	}()

	updates := []*model.UserEs{cloneUser(original[6]), cloneUser(original[8])}
	updates[0].Nickname = original[6].Nickname + "_updated"
	updates[0].Phone = "17800000000"
	updates[0].Identity = "worker"
	updates[1].Nickname = original[8].Nickname + "_updated"
	updates[1].Phone = "17800000001"
	updates[1].Identity = "student"

	if err := userES.BatchUpdate(ctx, updates); err != nil {
		t.Fatalf("BatchUpdate 执行失败: %v", err)
	}

	for _, expect := range updates {
		got := fetchUser(ctx, t, client, index, expect.ID)
		if got.Nickname != expect.Nickname {
			t.Fatalf("ID=%d 昵称未更新，期望=%s，得到=%s", expect.ID, expect.Nickname, got.Nickname)
		}
		if got.Phone != expect.Phone {
			t.Fatalf("ID=%d 手机号未更新，期望=%s，得到=%s", expect.ID, expect.Phone, got.Phone)
		}
		if got.Identity != expect.Identity {
			t.Fatalf("ID=%d 身份未更新，期望=%s，得到=%s", expect.ID, expect.Identity, got.Identity)
		}
	}
}

// newESTestClient 辅助函数1-新建ES连接
func newESTestClient(t *testing.T, cfg *conf.Config) *elastic.Client {
	t.Helper()
	url := fmt.Sprintf("http://%s:%d", cfg.Elastic.Host, cfg.Elastic.Port)
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		t.Fatalf("创建 ES 客户端失败: %v", err)
	}
	return client
}

// fetchUser 辅助函数2-获取特定索引中特定id的文档
func fetchUser(ctx context.Context, t *testing.T, client *elastic.Client, index string, id uint64) *model.UserEs {
	t.Helper()
	resp, err := client.Get().Index(index).Id(strconv.FormatUint(id, 10)).Do(ctx)
	if err != nil {
		t.Fatalf("获取文档 %d 失败: %v", id, err)
	}
	var user model.UserEs
	if err := json.Unmarshal(resp.Source, &user); err != nil {
		t.Fatalf("反序列化文档 %d 失败: %v", id, err)
	}
	return &user
}

// cloneUser 辅助函数3-深拷贝
func cloneUser(src *model.UserEs) *model.UserEs {
	if src == nil {
		return nil
	}
	clone := *src
	return &clone
}
