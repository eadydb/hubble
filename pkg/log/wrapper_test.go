package log

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantL   Level
		wantErr bool
	}{
		{
			name: "debug",
			args: args{
				s: "debug",
			},
			wantL: LevelDebug,
		},
		{
			name: "info",
			args: args{
				s: "info",
			},
			wantL: LevelInfo,
		},
		{
			name: "-4",
			args: args{
				s: "-4",
			},
			wantL: LevelDebug,
		},
		{
			name: "0",
			args: args{
				s: "0",
			},
			wantL: LevelInfo,
		},
		{
			name: "4",
			args: args{
				s: "4",
			},
			wantL: LevelWarn,
		},
		{
			name: "8",
			args: args{
				s: "8",
			},
			wantL: LevelError,
		},
		{
			name: "info+1",
			args: args{
				s: "info+1",
			},
			wantL: LevelInfo + 1,
		},
		{
			name: "info-1",
			args: args{
				s: "info-1",
			},
			wantL: LevelInfo - 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL, err := ParseLevel(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotL != tt.wantL {
				t.Errorf("ParseLevel() gotL = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}
