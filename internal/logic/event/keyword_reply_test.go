package event

import "testing"

func Test_decreasePlaceholderIndex(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Single placeholder with number",
			args: args{text: "{placeholder1}"},
			want: "{placeholder}",
		},
		{
			name: "Single placeholder without number",
			args: args{text: "{placeholder}"},
			want: "{placeholder}",
		},
		{
			name: "Multiple placeholders with numbers",
			args: args{text: "{placeholder1} and {placeholder2}"},
			want: "{placeholder} and {placeholder1}",
		},
		{
			name: "Placeholder with zero",
			args: args{text: "{placeholder0}"},
			want: "{placeholder}",
		},
		{
			name: "Placeholder with non-number suffix",
			args: args{text: "{placeholderA}"},
			want: "{placeholderA}",
		},
		{
			name: "Mixed placeholders",
			args: args{text: "{placeholder1} and {another2}"},
			want: "{placeholder} and {another1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decreasePlaceholderIndex(tt.args.text); got != tt.want {
				t.Errorf("decreasePlaceholderIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
