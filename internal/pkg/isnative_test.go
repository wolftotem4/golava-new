package pkg

import "testing"

func TestIsNativePackage(t *testing.T) {
	tests := []struct {
		name string
		pkg  string
		want bool
	}{
		{
			name: "native package",
			pkg:  "context",
			want: true,
		},
		{
			name: "non-native package",
			pkg:  "github.com/wolftotem4/golava-core/cookie",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNativePackage(tt.pkg); got != tt.want {
				t.Errorf("IsNativePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
