package spacetraders

import "testing"

func Test_isValidPlatformChannel(t *testing.T) {
	type args struct {
		platform string
		channel  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"discord", args{"discord", "wrong"}, false},
		{"general", args{"discord", "811472319876562989"}, false},
		{"sandbox", args{"discord", "846297823624298517"}, true},
		{"spacetraders", args{"discord", "1127471366501834763"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPlatformChannel(tt.args.platform, tt.args.channel); got != tt.want {
				t.Errorf("isValidPlatformChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
