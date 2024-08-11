package scraping

import "testing"

func TestNormalizePrice(t *testing.T) {
	type args struct {
		price string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "sucessfulty parse a BR price",
			args: args{
				price: "R$ 4.080,90",
			},
			want: 4080.90,
		},
		{
			name: "sucessfulty parse a BR price with noise",
			args: args{
				price: "R$ 4.080,90R$0000,00",
			},
			want: 4080.90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePrice(tt.args.price); got != tt.want {
				t.Errorf("NormalizePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
