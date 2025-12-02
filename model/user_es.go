package model

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
