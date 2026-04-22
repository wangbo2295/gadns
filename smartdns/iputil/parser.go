package iputil

import (
	"net"
	"strings"
)

// ParseIPs parses a list of IP inputs which can be:
// - Single IPs (e.g., "1.1.1.1")
// - IP ranges (e.g., "1.1.1.1-1.1.1.10")
// - CIDR notation (e.g., "192.168.1.0/24")
// Returns a flattened list of all IP addresses
func ParseIPs(inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return []string{}, nil
	}

	var result []string

	for _, input := range inputs {
		// Validate input first
		if err := ValidateIPInput(input); err != nil {
			return nil, err
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Determine the format and parse accordingly
		if strings.Contains(input, "/") {
			// CIDR notation
			ips, err := parseCIDR(input)
			if err != nil {
				return nil, err
			}
			result = append(result, ips...)
		} else if strings.Contains(input, "-") {
			// IP range
			ips, err := parseRange(input)
			if err != nil {
				return nil, err
			}
			result = append(result, ips...)
		} else {
			// Single IP
			result = append(result, input)
		}
	}

	return result, nil
}

// parseCIDR expands a CIDR notation into a list of IP addresses
func parseCIDR(cidr string) ([]string, error) {
	// Parse CIDR
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, ErrInvalidCIDRFormat
	}

	var ips []string

	// Get the network IP
	networkIP := ipnet.IP

	// For IPv4
	if networkIP.To4() != nil {
		// Convert to 4-byte representation
		networkIP = networkIP.To4()

		// Calculate the broadcast address (last IP in range)
		ones, bits := ipnet.Mask.Size()
		numIPs := 1 << (bits - ones)

		// Start from the network address
		currentIP := make(net.IP, 4)
		copy(currentIP, networkIP)

		// Iterate through all IPs in the range
		for i := 0; i < numIPs; i++ {
			ips = append(ips, currentIP.String())
			currentIP = inc(currentIP)
		}
	} else {
		// For IPv6
		ones, bits := ipnet.Mask.Size()
		numIPs := 1 << (bits - ones)

		// Limit the number of IPs to prevent memory issues
		if numIPs > 1000000 {
			numIPs = 1000000
		}

		// Start from the network address
		currentIP := make(net.IP, 16)
		copy(currentIP, networkIP)

		// Iterate through all IPs in the range
		for i := 0; i < numIPs; i++ {
			ips = append(ips, currentIP.String())
			currentIP = inc(currentIP)
		}
	}

	return ips, nil
}

// parseRange expands an IP range (e.g., "1.1.1.1-1.1.1.10") into a list of IP addresses
func parseRange(input string) ([]string, error) {
	// Validate the range format first
	if err := ValidateIPRange(input); err != nil {
		return nil, err
	}

	// Split by dash
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return nil, ErrInvalidIPRangeFormat
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	endIP := net.ParseIP(strings.TrimSpace(parts[1]))

	if startIP == nil || endIP == nil {
		return nil, ErrInvalidIPFormat
	}

	// Convert to consistent format (IPv4 or IPv6)
	startIP = startIP.To4()
	endIP = endIP.To4()

	if startIP == nil || endIP == nil {
		// IPv6 handling
		startIP = net.ParseIP(strings.TrimSpace(parts[0])).To16()
		endIP = net.ParseIP(strings.TrimSpace(parts[1])).To16()
		if startIP == nil || endIP == nil {
			return nil, ErrInvalidIPFormat
		}
	}

	var ips []string

	// Start from the beginning IP
	currentIP := make(net.IP, len(startIP))
	copy(currentIP, startIP)

	// Iterate until we reach the end IP
	for ipCompare(currentIP, endIP) <= 0 {
		ips = append(ips, currentIP.String())
		currentIP = inc(currentIP)
	}

	return ips, nil
}

// inc increments an IP address by 1
// Returns a new IP with the incremented value
func inc(ip net.IP) net.IP {
	// Make a copy to avoid modifying the original
	ip = make(net.IP, len(ip))
	copy(ip, ip)

	// Increment from right to left (network byte order)
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			// No overflow, we're done
			break
		}
		// Overflow, continue to next byte
	}

	return ip
}

// ipCompare compares two IP addresses byte by byte
// Returns -1 if a < b, 0 if a == b, 1 if a > b
func ipCompare(a, b []byte) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}

	// If all compared bytes are equal, check lengths
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}

	return 0
}
