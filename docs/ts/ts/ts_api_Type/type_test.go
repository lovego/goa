package ts_api_Type

import "testing"

func TestGenTsApiType(t *testing.T) {
	type args struct {
		outPath string
		apiFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// gt tpl -p=./test/accounts.yaml -o=./sql/erp/accounts.typings.d.ts
		{
			name: "tstype",
			args: args{
				outPath: "./test/accounts.yaml",
				apiFile: "/Users/lchjczw/work/project/erp/erp-api/api/docment/doc/account/account.api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenApiType(tt.args.outPath, tt.args.apiFile); (err != nil) != tt.wantErr {
				t.Errorf("GenApiType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
