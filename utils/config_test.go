package utils

import "testing"

func TestConfigOrDefault(t *testing.T) {
    tests := []struct {
        name     string
        key      string
        fallback string
        want     string
    }{
        {"missing key uses fallback", "nonexistent_key", "http://fallback.com", "http://fallback.com"},
        {"empty fallback", "nonexistent_key", "", ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := configOrDefault(tt.key, tt.fallback)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}