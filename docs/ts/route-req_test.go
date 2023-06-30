package ts

import "testing"

func Test_getTsParamPath(t *testing.T) {
	type args struct {
		fullPath string
		names    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "param",
			args: args{
				fullPath: `/123/(?P<type>(QD|KD))/(?P<id>\d+)/(?P<action>(submit|revoke|pass|costPass|return|refuse|financial))/wer`,
				names:    []string{"type", "id", "action"},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTsParamPath(tt.args.fullPath, tt.args.names); got != tt.want {
				t.Errorf("getTsParamPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
