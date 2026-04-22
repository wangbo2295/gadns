package iputil

import (
	"net"
	"reflect"
	"testing"
)

func TestParseIPs(t *testing.T) {
	tests := []struct {
		name      string
		inputs    []string
		want      []string
		wantErr   bool
		errChecks []error // Optional: check for specific errors
	}{
		{
			name:    "single IP",
			inputs:  []string{"1.1.1.1"},
			want:    []string{"1.1.1.1"},
			wantErr: false,
		},
		{
			name:    "multiple IPs",
			inputs:  []string{"1.1.1.1", "2.2.2.2"},
			want:    []string{"1.1.1.1", "2.2.2.2"},
			wantErr: false,
		},
		{
			name:    "IP range small",
			inputs:  []string{"1.1.1.1-1.1.1.3"},
			want:    []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"},
			wantErr: false,
		},
		{
			name:    "IP range single",
			inputs:  []string{"192.168.1.1-192.168.1.1"},
			want:    []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "IP range larger",
			inputs:  []string{"10.0.0.1-10.0.0.5"},
			want:    []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4", "10.0.0.5"},
			wantErr: false,
		},
		{
			name:    "CIDR /30",
			inputs:  []string{"1.1.1.0/30"},
			want:    []string{"1.1.1.0", "1.1.1.1", "1.1.1.2", "1.1.1.3"},
			wantErr: false,
		},
		{
			name:    "CIDR /32",
			inputs:  []string{"192.168.1.1/32"},
			want:    []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "CIDR /24",
			inputs:  []string{"192.168.1.0/24"},
			wantErr: false, // We'll check length separately
		},
		{
			name:    "mixed formats",
			inputs:  []string{"1.1.1.1", "1.1.1.5-1.1.1.7", "2.2.2.0/30"},
			want:    []string{"1.1.1.1", "1.1.1.5", "1.1.1.6", "1.1.1.7", "2.2.2.0", "2.2.2.1", "2.2.2.2", "2.2.2.3"},
			wantErr: false,
		},
		{
			name:    "mixed with multiple single IPs",
			inputs:  []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			want:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			wantErr: false,
		},
		{
			name:    "empty input list",
			inputs:  []string{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "nil input list",
			inputs:  nil,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "empty string in list",
			inputs:  []string{""},
			wantErr: true,
			errChecks: []error{ErrEmptyInput},
		},
		{
			name:    "invalid IP format",
			inputs:  []string{"invalid"},
			wantErr: true,
			errChecks: []error{ErrInvalidIPFormat},
		},
		{
			name:    "invalid range format",
			inputs:  []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "invalid CIDR format",
			inputs:  []string{"192.168.1.0/33"},
			wantErr: true,
		},
		{
			name:    "reversed IP range",
			inputs:  []string{"192.168.1.100-192.168.1.1"},
			wantErr: true,
			errChecks: []error{ErrInvalidRange},
		},
		{
			name:    "invalid IP in range",
			inputs:  []string{"invalid-192.168.1.100"},
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid",
			inputs:  []string{"1.1.1.1", "invalid", "2.2.2.2"},
			wantErr: true,
		},
		{
			name:    "whitespace handling",
			inputs:  []string{" 1.1.1.1 ", " 2.2.2.2 "},
			want:    []string{"1.1.1.1", "2.2.2.2"},
			wantErr: false,
		},
		{
			name:    "IPv6 single",
			inputs:  []string{"::1"},
			want:    []string{"::1"},
			wantErr: false,
		},
		{
			name:    "IPv6 single", // Changed from CIDR to single IP to avoid large address space
			inputs:  []string{"2001:db8::1"},
			want:    []string{"2001:db8::1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPs(tt.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIPs(%v) error = %v, wantErr %v", tt.inputs, err, tt.wantErr)
				return
			}

			// Check for specific error types if provided
			if tt.wantErr && len(tt.errChecks) > 0 {
				errMatched := false
				for _, expectedErr := range tt.errChecks {
					if err == expectedErr {
						errMatched = true
						break
					}
				}
				if !errMatched {
					t.Errorf("ParseIPs(%v) error = %v, wanted one of %v", tt.inputs, err, tt.errChecks)
				}
				return
			}

			// If we want exact match and no error
			if !tt.wantErr && len(tt.want) > 0 && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIPs(%v) = %v, want %v", tt.inputs, got, tt.want)
			}

			// Special case for /24 - just check length
			if tt.name == "CIDR /24" && err == nil {
				if len(got) != 256 {
					t.Errorf("ParseIPs(%v) returned %d IPs, want 256", tt.inputs, len(got))
				}
			}
		})
	}
}

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantLen int // Expected number of IPs (use -1 to skip length check)
		wantErr bool
	}{
		{
			name:    "IPv4 /32",
			cidr:    "192.168.1.1/32",
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "IPv4 /30",
			cidr:    "1.1.1.0/30",
			wantLen: 4,
			wantErr: false,
		},
		{
			name:    "IPv4 /24",
			cidr:    "192.168.1.0/24",
			wantLen: 256,
			wantErr: false,
		},
		{
			name:    "IPv4 /22", // Smaller range for faster tests
			cidr:    "10.0.0.0/22",
			wantLen: 1024,
			wantErr: false,
		},
		{
			name:    "invalid CIDR",
			cidr:    "invalid",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			cidr:    "192.168.1.0/33",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCIDR(%q) error = %v, wantErr %v", tt.cidr, err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.wantLen >= 0 && len(got) != tt.wantLen {
				t.Errorf("parseCIDR(%q) returned %d IPs, want %d", tt.cidr, len(got), tt.wantLen)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "small range",
			input:   "1.1.1.1-1.1.1.3",
			want:    []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"},
			wantErr: false,
		},
		{
			name:    "single IP range",
			input:   "192.168.1.1-192.168.1.1",
			want:    []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "consecutive IPs",
			input:   "10.0.0.1-10.0.0.2",
			want:    []string{"10.0.0.1", "10.0.0.2"},
			wantErr: false,
		},
		{
			name:    "invalid format no dash",
			input:   "192.168.1.1",
			wantErr: true,
		},
		{
			name:    "invalid format multiple dashes",
			input:   "192.168.1.1-192.168.1.2-192.168.1.3",
			wantErr: true,
		},
		{
			name:    "invalid start IP",
			input:   "invalid-192.168.1.2",
			wantErr: true,
		},
		{
			name:    "invalid end IP",
			input:   "192.168.1.1-invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRange(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInc(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{"increment last octet", "1.1.1.1", "1.1.1.2"},
		{"increment to next octet", "1.1.1.255", "1.1.2.0"},
		{"increment boundary", "1.1.255.255", "1.2.0.0"},
		{"increment max IPv4 wraps to zero", "255.255.255.255", "0.0.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			got := inc(ip).String()
			if got != tt.want {
				t.Errorf("inc(%q) = %q, want %q", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIPCompare(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int // -1 if a < b, 0 if a == b, 1 if a > b
	}{
		{"a less than b", "1.1.1.1", "1.1.1.2", -1},
		{"a equal b", "1.1.1.1", "1.1.1.1", 0},
		{"a greater than b", "1.1.1.2", "1.1.1.1", 1},
		{"different octets", "192.168.1.1", "192.168.2.1", -1},
		{"different subnet", "10.0.0.1", "192.168.1.1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := net.ParseIP(tt.a)
			b := net.ParseIP(tt.b)
			if a == nil || b == nil {
				t.Fatalf("Failed to parse IPs: %s, %s", tt.a, tt.b)
			}
			got := ipCompare(a.To4(), b.To4())
			if got != tt.want {
				t.Errorf("ipCompare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
