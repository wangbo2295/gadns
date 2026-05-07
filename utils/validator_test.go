package utils

import (
	"testing"
)

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.1", false},
		{"valid IPv4 localhost", "127.0.0.1", false},
		{"valid IPv4", "8.8.8.8", false},
		{"valid IPv4", "255.255.255.255", false},
		{"empty string", "", true},
		{"hostname", "example.com", true},
		{"not an IP", "not-an-ip", true},
		{"out of range", "256.1.1.1", true},
		{"IPv6", "2001:db8::1", true},
		{"IPv6 localhost", "::1", true},
		{"CIDR notation", "192.168.1.0/24", true},
		{"IP range", "1.1.1.1-1.1.1.10", true},
		{"special chars", "!@#$%", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIP(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIP(%q) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
			}
		})
	}
}
