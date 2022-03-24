package ynot

import "testing"

func Test_ynot(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"why don't use ruby?"}, true},
		{"2", args{"use ruby why don't not?"}, false},
		{"3", args{"I don't see why not java. RIght?"}, false},
		{"4", args{"why don't you try amazon?1"}, true},
		{"5", args{"y not just use mysql?"}, true},
		{"6", args{"Why don't we just pay in bitcoin instead?"}, true},
		{"cobol1", args{"why don't we just code in COBOL"}, true},  // Reported by SM
		{"cobol2", args{"why don't we just code in COBOL?"}, true}, // Reported by SM
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ynot(tt.args.i); got != tt.want {
				t.Errorf("ynot() = %v, want %v", got, tt.want)
			}
		})
	}
}
