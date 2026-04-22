package smartdns

import (
	"testing"
)

func TestSmartDNSInterface(t *testing.T) {
	var _ SmartDNS = (*mockProvider)(nil)
}

type mockProvider struct{}

func (m *mockProvider) Create(name string, ips []string) (*Record, error) {
	return &Record{Name: name, CNAME: name + ".test.com", IPs: ips}, nil
}

func (m *mockProvider) Update(name string, ips []string) (*Record, error) {
	return &Record{Name: name, CNAME: name + ".test.com", IPs: ips}, nil
}

func (m *mockProvider) Get(name string) (*Record, error) {
	return &Record{Name: name, CNAME: name + ".test.com", IPs: []string{}}, nil
}

func (m *mockProvider) List() ([]*Record, error) {
	return []*Record{}, nil
}

func (m *mockProvider) Delete(name string) error {
	return nil
}
