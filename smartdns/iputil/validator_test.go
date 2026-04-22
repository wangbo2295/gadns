package iputil

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
		{"valid IPv4 zero", "0.0.0.0", false},
		{"valid IPv6", "2001:db8::1", false},
		{"valid IPv6 localhost", "::1", false},
		{"valid IPv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"empty string", "", true},
		{"not an IP", "not-an-ip", true},
		{"invalid IPv4 - missing octet", "192.168.1", true},
		{"invalid IPv4 - extra octet", "192.168.1.1.1", true},
		{"invalid IPv4 - out of range", "256.1.1.1", true},
		{"invalid IPv4 - negative", "-1.1.1.1", true},
		{"invalid IPv4 - text", "192.168.1.a", true},
		{"invalid IPv6", "2001:db8:::1", true},
		{"valid IPv6 documentation prefix", "2001:db8::", false},
		{"hostname", "example.com", true},
		{"hostname with subdomain", "www.example.com", true},
		{"just text", "hello", true},
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

func TestValidateIPRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid range", "192.168.1.1-192.168.1.100", false},
		{"valid range same first three octets", "10.0.0.1-10.0.0.255", false},
		{"valid range different subnets", "192.168.1.100-192.168.2.50", false},
		{"valid range single IP", "192.168.1.1-192.168.1.1", false},
		{"valid range sequential", "1.2.3.4-1.2.3.5", false},
		{"valid range large span", "10.0.0.1-10.0.1.1", false},
		{"empty string", "", true},
		{"no dash", "192.168.1.1 192.168.1.100", true},
		{"missing start IP", "-192.168.1.100", true},
		{"missing end IP", "192.168.1.1-", true},
		{"multiple dashes", "192.168.1.1-192.168.1.100-192.168.1.200", true},
		{"invalid start IP", "invalid-192.168.1.100", true},
		{"invalid end IP", "192.168.1.1-invalid", true},
		{"both invalid IPs", "invalid-invalid", true},
		{"out of range start", "256.1.1.1-192.168.1.100", true},
		{"out of range end", "192.168.1.1-256.1.1.1", true},
		{"reversed range (start > end)", "192.168.1.100-192.168.1.1", true},
		{"reversed range different subnet", "192.168.2.50-192.168.1.100", true},
		{"start octet greater than end", "192.168.1.5-192.168.1.2", true},
		{"just IP no range", "192.168.1.1", true},
		{"CIDR notation", "192.168.1.0/24", true},
		{"hostname range", "example.com-example.org", true},
		{"just text", "hello-world", true},
		{"dash but not IPs", "a-b", true},
		{"spaces around", " 192.168.1.1-192.168.1.100 ", false},
		{"tabs around", "\t192.168.1.1-192.168.1.100\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCIDR(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid CIDR /24", "192.168.1.0/24", false},
		{"valid CIDR /32", "192.168.1.1/32", false},
		{"valid CIDR /16", "10.0.0.0/16", false},
		{"valid CIDR /8", "10.0.0.0/8", false},
		{"valid CIDR /0", "0.0.0.0/0", false},
		{"valid single IP", "192.168.1.1", false}, // Should accept single IP as valid
		{"valid localhost", "127.0.0.1", false},   // Single IP
		{"empty string", "", true},
		{"missing prefix", "192.168.1.0/", true},
		{"missing IP", "/24", true},
		{"invalid IP format", "192.168.1/24", true},
		{"invalid prefix too large", "192.168.1.0/33", true},
		{"invalid prefix negative", "192.168.1.0/-1", true},
		{"invalid prefix text", "192.168.1.0/abc", true},
		{"multiple slashes", "192.168.1.0/24/16", true},
		{"IP range format", "192.168.1.1-192.168.1.100", true},
		{"hostname", "example.com", true},
		{"hostname with CIDR", "example.com/24", true},
		{"just text", "hello", true},
		{"spaces around", " 192.168.1.0/24 ", false},
		{"invalid IP out of range", "256.1.1.0/24", true},
		{"prefix zero no slash", "192.168.1.0", false}, // Should be valid as single IP
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCIDR(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCIDR(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateIPInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Single IP tests
		{"single IPv4", "192.168.1.1", false},
		{"single IPv6", "2001:db8::1", false},
		{"single localhost", "127.0.0.1", false},

		// IP range tests
		{"IP range valid", "192.168.1.1-192.168.1.100", false},
		{"IP range single", "192.168.1.1-192.168.1.1", false},

		// CIDR tests
		{"CIDR valid", "192.168.1.0/24", false},
		{"CIDR /32", "192.168.1.1/32", false},

		// Invalid inputs
		{"empty string", "", true},
		{"hostname", "example.com", true},
		{"just text", "hello", true},
		{"invalid IP", "not-an-ip", true},
		{"invalid range", "192.168.1.100-192.168.1.1", true},
		{"invalid CIDR", "192.168.1.0/33", true},
		{"special chars", "!@#$%", true},
		{"mixed format", "192.168.1.1/24-192.168.1.100", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPInput(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
