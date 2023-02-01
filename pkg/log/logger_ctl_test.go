package log

import (
	"errors"
	"testing"

	"golang.org/x/exp/slog"
)

func Test_formatValue(t *testing.T) {
	type args struct {
		val slog.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "format for error",
			args: args{
				val: slog.AnyValue(errors.New("unknown command \"subcommand\" for \"kwokctl\"")),
			},
			want: quoteIfNeed(errors.New("unknown command \"subcommand\" for \"kwokctl\"").Error()),
		},
		{
			name: "format for string",
			args: args{
				val: slog.AnyValue("unknown command \"subcommand\" for \"kwokctl\""),
			},
			want: quoteIfNeed("unknown command \"subcommand\" for \"kwokctl\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatValue(tt.args.val); got != tt.want {
				t.Errorf("formatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
