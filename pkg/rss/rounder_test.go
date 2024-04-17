package rss

import (
	"testing"
)

func TestRound(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Тест1",
			args: args{url: "https://cprss.s3.amazonaws.com/golangweekly.com.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Round(tt.args.url)
			for _, item := range got.Channel.Item {
				t.Logf("Title: %s", item.Title)
			}

		})
	}
}