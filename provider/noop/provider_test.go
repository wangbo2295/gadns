package noop_test

import (
	"testing"

	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/provider/noop"
)

func TestNoopImplementsInterface(t *testing.T) {
	provider := noop.NewProvider(&noop.Config{Domain: "example.com"})
	var _ core.CNAMEProvider = provider
}

func TestNoopCreateAndGet(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	r, err := p.Create("app.example.com", []string{"1.1.1.1"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if r.CNAME == "" || r.CNAME == r.Name {
		t.Errorf("CNAME = %v, expected hashed CNAME", r.CNAME)
	}

	got, err := p.Get("app.example.com")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.IPs) != 1 || got.IPs[0] != "1.1.1.1" {
		t.Errorf("Get() IPs = %v", got.IPs)
	}
}

func TestNoopCreateDuplicate(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	_, err := p.Create("app.example.com", []string{"1.1.1.1"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Create("app.example.com", []string{"2.2.2.2"})
	if err == nil {
		t.Error("expected error for duplicate record")
	}
}

func TestNoopUpdate(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	p.Create("app.example.com", []string{"1.1.1.1"})
	r, err := p.Update("app.example.com", []string{"2.2.2.2"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if r.IPs[0] != "2.2.2.2" {
		t.Errorf("Update() IPs = %v", r.IPs)
	}
}

func TestNoopUpdateNotFound(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	_, err := p.Update("nope.example.com", []string{"1.1.1.1"})
	if err == nil {
		t.Error("expected error for non-existent record")
	}
}

func TestNoopListAndDelete(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	p.Create("a.example.com", []string{"1.1.1.1"})
	p.Create("b.example.com", []string{"2.2.2.2"})

	list, err := p.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("List() len = %d, want 2", len(list))
	}

	if err := p.Delete("a.example.com"); err != nil {
		t.Fatal(err)
	}

	list, _ = p.List()
	if len(list) != 1 {
		t.Errorf("List() after delete len = %d, want 1", len(list))
	}
}

func TestNoopGetNotFound(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	_, err := p.Get("nope.example.com")
	if err == nil {
		t.Error("expected error for non-existent record")
	}
}

func TestNoopDeleteNotFound(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	err := p.Delete("nope.example.com")
	if err == nil {
		t.Error("expected error for non-existent record")
	}
}

func TestNoopInvalidIP(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	_, err := p.Create("app.example.com", []string{"not-an-ip"})
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}

func TestNoopIPv6Rejected(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	_, err := p.Create("app.example.com", []string{"2001:db8::1"})
	if err == nil {
		t.Error("expected error for IPv6")
	}
}

func TestNoopMultipleIPs(t *testing.T) {
	p := noop.NewProvider(&noop.Config{Domain: "example.com"})

	r, err := p.Create("web.example.com", []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.IPs) != 3 {
		t.Errorf("Create() IPs len = %d, want 3", len(r.IPs))
	}
}
