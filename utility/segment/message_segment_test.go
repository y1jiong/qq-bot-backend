package segment

import (
	"reflect"
	"testing"
)

func BenchmarkParseMessage(b *testing.B) {
	message := "1[CQ:custom,data=123,flag=1]2[CQ:custom,data=456][CQ:custom,data=789]3"

	b.ReportAllocs()
	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseMessage(message)
	}
}

func TestParseMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want messageSegments
	}{
		{
			name: "normal",
			args: args{
				message: "1[CQ:custom,data=123,flag=1]2[CQ:custom,data=456][CQ:custom,data=789]3",
			},
			want: messageSegments{
				{
					Type: "text",
					Data: map[string]string{
						"text": "1",
					},
				},
				{
					Type: "custom",
					Data: map[string]string{
						"data": "123",
						"flag": "1",
					},
				},
				{
					Type: "text",
					Data: map[string]string{
						"text": "2",
					},
				},
				{
					Type: "custom",
					Data: map[string]string{
						"data": "456",
					},
				},
				{
					Type: "custom",
					Data: map[string]string{
						"data": "789",
					},
				},
				{
					Type: "text",
					Data: map[string]string{
						"text": "3",
					},
				},
			},
		},
		{
			name: "escape",
			args: args{
				message: "[CQ:custom,&#44;=&#91;&amp;&#93;]",
			},
			want: messageSegments{
				{
					Type: "custom",
					Data: map[string]string{
						"&#44;": "[&]",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMessage(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
