package smartdns

type Record struct {
	Name   string
	CNAME  string
	IPs    []string
}
