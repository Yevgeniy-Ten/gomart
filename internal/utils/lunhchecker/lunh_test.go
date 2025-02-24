package lunhchecker

import "testing"

func TestLuhnCheck(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid",
			number: "79927398713",
			want:   true,
		},
		{
			name:   "invalid",
			number: "79927398710",
			want:   false,
		},
		{
			name:   "invalid char",
			number: "7992739871a",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnCheck(tt.number); got != tt.want {
				t.Errorf("LuhnCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
