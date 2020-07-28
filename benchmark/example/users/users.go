package users

import "time"

type ListReq struct {
	Name     string `c:"用户名称"`
	Type     string `c:"用户类型"`
	Page     int    `c:"页码"`
	PageSize int    `c:"每页数据条数"`
}

type ListResp struct {
	TotalSize int `c:"总数据条数"`
	TotalPage int `c:"总页数"`
	Rows      []struct {
		Id    int    `c:"ID"`
		Name  string `c:"名称"`
		Phone string `c:"电话号码"`
	}
}

type Session struct {
	UserId  int64
	LoginAt time.Time
}

func (l *ListReq) Run(sess *Session) (ListResp, error) {
	return ListResp{}, nil
}

type DetailResp struct {
	TotalSize int `c:"总数据条数"`
	TotalPage int `c:"总页数"`
	Rows      []struct {
		Id    int    `c:"ID"`
		Name  string `c:"名称"`
		Phone string `c:"电话号码"`
	}
}

func Detail(userId int64) (DetailResp, error) {
	return DetailResp{}, nil
}
