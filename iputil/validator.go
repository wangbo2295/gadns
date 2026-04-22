package iputil

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrEmptyInput           = errors.New("empty input")
	ErrInvalidIPFormat      = errors.New("invalid IP format")
	ErrInvalidIPRangeFormat = errors.New("invalid IP range format, must be x.x.x.x-x.x.x.x")
	ErrInvalidCIDRFormat    = errors.New("invalid CIDR format")
	ErrInvalidRange         = errors.New("start IP must be less than or equal to end IP")
)

// ValidateIP validates if the input is a valid IP address (IPv4 or IPv6)
func ValidateIP(ip string) error {
	if ip == "" {
		return ErrEmptyInput
	}

	// Trim whitespace
	ip = strings.TrimSpace(ip)

	// net.ParseIP accepts both IPv4 and IPv6
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ErrInvalidIPFormat
	}

	return nil
}

// ValidateIPRange validates if the input is a valid IP range in format x.x.x.x-x.x.x.x
// This is a strict format validator - it requires exactly two IPs separated by a dash
func ValidateIPRange(input string) error {
	if input == "" {
		return ErrEmptyInput
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Split by dash
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return ErrInvalidIPRangeFormat
	}

	startIPStr := strings.TrimSpace(parts[0])
	endIPStr := strings.TrimSpace(parts[1])

	// Check for empty parts
	if startIPStr == "" || endIPStr == "" {
		return ErrInvalidIPRangeFormat
	}

	// Parse start IP
	startIP := net.ParseIP(startIPStr)
	if startIP == nil {
		return ErrInvalidIPFormat
	}

	// Parse end IP
	endIP := net.ParseIP(endIPStr)
	if endIP == nil {
		return ErrInvalidIPFormat
	}

	// Compare IPs byte by byte
	// net.ParseIP returns IPv4 in IPv6-mapped format, so we need to handle that
	startIP = startIP.To4()
	endIP = endIP.To4()

	// If both are IPv4, compare as 4-byte slices
	if startIP != nil && endIP != nil {
		for i := 0; i < 4; i++ {
			if startIP[i] < endIP[i] {
				return nil // Valid: start < end
			}
			if startIP[i] > endIP[i] {
				return ErrInvalidRange // Invalid: start > end
			}
			// If equal, continue to next byte
		}
		return nil // All bytes equal, valid (single IP)
	}

	// For IPv6, compare as 16-byte slices
	startIP16 := startIP.To16()
	endIP16 := endIP.To16()

	if startIP16 == nil || endIP16 == nil {
		return ErrInvalidIPFormat
	}

	for i := 0; i < 16; i++ {
		if startIP16[i] < endIP16[i] {
			return nil // Valid: start < end
		}
		if startIP16[i] > endIP16[i] {
			return ErrInvalidRange // Invalid: start > end
		}
		// If equal, continue to next byte
	}

	return nil // All bytes equal, valid (single IP)
}

// ValidateCIDR validates if the input is a valid CIDR notation or single IP
// Accepts both CIDR notation (e.g., 192.168.1.0/24) and single IPs (e.g., 192.168.1.1)
func ValidateCIDR(input string) error {
	if input == "" {
		return ErrEmptyInput
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Check if it contains a slash (CIDR notation)
	if strings.Contains(input, "/") {
		// Try to parse as CIDR
		_, _, err := net.ParseCIDR(input)
		if err != nil {
			return ErrInvalidCIDRFormat
		}
		return nil
	}

	// If no slash, validate as single IP
	return ValidateIP(input)
}

// ValidateIPInput validates any IP input format
// It can be: single IP, IP range (x.x.x.x-x.x.x.x), or CIDR notation
func ValidateIPInput(input string) error {
	if input == "" {
		return ErrEmptyInput
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Check if it's a CIDR notation
	if strings.Contains(input, "/") {
		return ValidateCIDR(input)
	}

	// Check if it's an IP range
	if strings.Contains(input, "-") {
		return ValidateIPRange(input)
	}

	// Otherwise, treat as single IP
	return ValidateIP(input)
}
