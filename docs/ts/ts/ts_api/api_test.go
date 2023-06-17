package ts_api

import (
	"testing"
)

func TestGenTsApi(t *testing.T) {
	type args struct {
		outPath string
		apiFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// gt tpl -p=./test/accounts.yaml -o=./sql/erp/accounts.ts
		{
			name: "ts",
			args: args{
				outPath: "./test/api.yaml",
				apiFile: "/Users/lchjczw/work/project/erp/erp-api/api/doc/accounts/accounts.api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenApi(tt.args.outPath, tt.args.apiFile); (err != nil) != tt.wantErr {
				t.Errorf("GenApi() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
