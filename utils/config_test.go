package utils

import (
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestConfigOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string // empty means don't set
		fallback string
		want     string
	}{
		{
			name:     "missing key returns fallback",
			key:      "nonexistent_key_xyz",
			fallback: "http://fallback.com",
			want:     "http://fallback.com",
		},
		{
			name:     "empty fallback on missing key returns empty",
			key:      "nonexistent_key_abc",
			fallback: "",
			want:     "",
		},
		{
			name:     "set key returns value not fallback",
			key:      "test_config_key",
			setValue: "http://real.com",
			fallback: "http://fallback.com",
			want:     "http://real.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" {
				beego.AppConfig.Set(tt.key, tt.setValue)
			}
			got := configOrDefault(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}