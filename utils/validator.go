package utils

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrEmptyInput      = errors.New("empty input")
	ErrInvalidIPFormat = errors.New("invalid IPv4 address")
)

// ValidateIP validates if the input is a valid IPv4 address
func ValidateIP(ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ErrEmptyInput
	}

	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() == nil {
		return ErrInvalidIPFormat
	}

	return nil
}
