package ts_tpl

import (
	"testing"

	"github.com/lovego/goa/docs/ts/ts/api_type"
)

func TestApiTypeInfo_Run(t *testing.T) {
	type fields struct {
		File     string
		Title    string
		Desc     string
		Method   string
		Router   string
		TypeList []api_type.Object
		Header   *api_type.Object
		Query    *api_type.Object
		Body     *api_type.Object
		Param    *api_type.Object
		Resp     *api_type.Object
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "tpl",
			fields: fields{
				File:     "./test.ts",
				Title:    "测试接口",
				Desc:     "测试接口",
				Method:   "post",
				Router:   "/api/user/update",
				TypeList: []api_type.Object{},
				Header:   nil,
				Query:    nil,
				Body: &api_type.Object{
					Name:     "User",
					JsonName: "user",
					Comment:  "用户信息",
					Members:  []api_type.Member{},
				},
				Param: nil,
				Resp: &api_type.Object{
					Name:     "User",
					JsonName: "user",
					Comment:  "用户信息",
					Members:  []api_type.Member{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ApiInfo{
				File:     tt.fields.File,
				Title:    tt.fields.Title,
				Desc:     tt.fields.Desc,
				Method:   tt.fields.Method,
				Router:   tt.fields.Router,
				TypeList: tt.fields.TypeList,
				Header:   tt.fields.Header,
				Query:    tt.fields.Query,
				Body:     tt.fields.Body,
				Param:    tt.fields.Param,
				Resp:     tt.fields.Resp,
			}
			if err := a.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
