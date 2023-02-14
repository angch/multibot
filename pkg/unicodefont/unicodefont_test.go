package unicodefont

import "testing"

func Test_unicodeReplace(t *testing.T) {
	type args struct {
		to int
		i  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"trad2simple", args{0, "橋頭"}, "桥头"},
		{"simple2trad", args{0, "桥头"}, "橋頭"},
		{"simple2trad", args{0, "新年快乐"}, "新年快樂"},
		{"trad2simple", args{0, "新年快樂"}, "新年快乐"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unicodeReplace(tt.args.to, tt.args.i); got != tt.want {
				t.Errorf("unicodeReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
