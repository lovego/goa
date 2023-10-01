package jet

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTpl(t *testing.T) {
	type args struct {
		tpl  []byte
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *bytes.Buffer
		wantErr bool
	}{
		// TODO: Add test cases.
		{

			//RunJetTest(t, data, nil, "actionNode_Field2", , `Oi José Santos<email@example.com>`)
			//
			name: "tpl",
			args: args{
				tpl: []byte(`<div> ≤.user.Name≥<≤.user.Email≥></div>`),
				data: map[string]map[string]string{
					"user": {
						"Name":  "lch",
						"Email": "lch@qq.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tpl(tt.args.tpl, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tpl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
