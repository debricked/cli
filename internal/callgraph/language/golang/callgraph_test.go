package golang

import (
	"testing"
)

func TestCleanSymbol(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "Test with pointer symbol",
			in:   "(*github.com/spf13/afero/mem.File).Open",
			want: "github.com/spf13/afero/mem.File.Open",
		},
		{
			name: "Test with non-pointer symbol",
			in:   "github.com/spf13/afero/mem.File.Open",
			want: "github.com/spf13/afero/mem.File.Open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanSymbol(tt.in); got != tt.want {
				t.Errorf("cleanSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
