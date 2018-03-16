package gelf

import "testing"

func TestZapLevelToGelfLevel(t *testing.T) {
	type args struct {
		l int32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "invalid level", args: args{l: -2}, want: 1},
		{name: "debug level", args: args{l: -1}, want: 7},
		{name: "info level", args: args{l: 0}, want: 6},
		{name: "warn level", args: args{l: 1}, want: 4},
		{name: "error level", args: args{l: 2}, want: 3},
		{name: "dpanic level", args: args{l: 3}, want: 0},
		{name: "panic level", args: args{l: 4}, want: 0},
		{name: "fatal level", args: args{l: 5}, want: 0},
		{name: "invalid level", args: args{l: 6}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZapLevelToGelfLevel(tt.args.l); got != tt.want {
				t.Errorf("ZapLevelToGelfLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// EOF
