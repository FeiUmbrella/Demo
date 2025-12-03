package model

import "github.com/olivere/elastic/v7"

type UserEs struct {
	ID         uint64 `json:"id,omitempty" mapstructure:"id"`
	Username   string `json:"username,omitempty" mapstructure:"username"`
	Nickname   string `json:"nickname,omitempty" mapstructure:"nickname"`
	Phone      string `json:"phone,omitempty" mapstructure:"phone"`
	Age        uint64 `json:"age,omitempty" mapstructure:"age"`
	Ancestral  string `json:"ancestral,omitempty" mapstructure:"Ancestral"`
	Identity   string `json:"identity,omitempty" mapstructure:"identity"`
	UpdateTime uint64 `json:"update_time,omitempty" mapstructure:"update_time"`
	CreateTime uint64 `json:"create_time,omitempty" mapstructure:"create_time"`
}

type SearchReq struct {
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Identity  string `json:"identity"`
	Ancestral string `json:"ancestral"`
	Num       int    `json:"num"`
	Size      int    `json:"size"`
}

// EsSearch 是 ES bool-query 的先决条件
type EsSearch struct {
	MustQuery    []elastic.Query
	MustNotQuery []elastic.Query
	ShouldQuery  []elastic.Query
	Filters      []elastic.Query
	Sorters      []elastic.Sorter
	From         int //分页
	Size         int //每页大小
}

// ToFilter 通过传入的req来构造ES的条件过滤
func (r *SearchReq) ToFilter() *EsSearch {
	var search EsSearch
	if r.Nickname != "" {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("nickname", r.Nickname))
	}
	if r.Phone != "" {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewTermsQuery("phone", r.Phone))
	}
	if r.Ancestral != "" {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("ancestral", r.Ancestral))
	}
	if r.Identity != "" {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("identity", r.Identity))
	}
	// 按创建时间降序排列
	if search.Sorters == nil {
		search.Sorters = append(search.Sorters, elastic.NewFieldSort("create_time").Desc())
	}

	search.From = (r.Num - 1) * r.Size
	search.Size = r.Size
	return &search
}
