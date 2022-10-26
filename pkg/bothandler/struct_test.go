package bothandler

import "testing"

func Test_sanitizeFilename(t *testing.T) {
	type args struct {
		f         string
		extension string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"codemonkey", args{"I don't care too much for money\nMoney can't buy me love", "jpg"}, "I_don't_care_too_much_for_money_Money_can't_buy_me_love.jpg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeFilename(tt.args.f, tt.args.extension); got != tt.want {
				t.Errorf("sanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
