package main

import "testing"

func Test_decode(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "12",
			args: args{
				message: "12",
			},
			want: 2,
		},
		{
			name: "226",
			args: args{
				message: "226",
			},
			want: 3,
		},
		{
			name: "06",
			args: args{
				message: "06",
			},
			want: 0,
		},
		{
			name: "0",
			args: args{
				message: "0",
			},
			want: 0,
		},
		{
			name: "106",
			args: args{
				message: "106",
			},
			want: 1,
		},
		{
			name: "1006",
			args: args{
				message: "1006",
			},
			want: 0,
		},
		{
			name: "2101",
			args: args{
				message: "2101",
			},
			want: 1,
		},
		{
			name: "2",
			args: args{
				message: "2",
			},
			want: 1,
		},
		{
			name: "22",
			args: args{
				message: "22",
			},
			want: 2,
		},
		{
			name: "221",
			args: args{
				message: "221",
			},
			want: 3,
		},
		{
			name: "2211",
			args: args{
				message: "2211",
			},
			want: 5,
		},
		{
			name: "22110",
			args: args{
				message: "22110",
			},
			want: 3,
		},
		{
			name: "221101",
			args: args{
				message: "221101",
			},
			want: 3,
		},
		{
			name: "2211011",
			args: args{
				message: "2211011",
			},
			want: 6,
		},
		{
			name: "230",
			args: args{
				message: "230",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decode(tt.args.message); got != tt.want {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
