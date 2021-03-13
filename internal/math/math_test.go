package math

import "testing"

func TestMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Max works as expected when first parameter is the maximum",
			args: args{
				a: 10,
				b: 1,
			},
			want: 10,
		},
		{
			name: "Max works as expected when second parameter is the maximum",
			args: args{
				a: 1,
				b: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Min works as expected when first parameter is the minimum",
			args: args{
				a: 10,
				b: 1,
			},
			want: 1,
		},
		{
			name: "Min works as expected when second parameter is the minimum",
			args: args{
				a: 1,
				b: 10,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}
