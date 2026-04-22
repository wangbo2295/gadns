package smartdns

type SmartDNS interface {
	Create(name string, ips []string) (*Record, error)
	Update(name string, ips []string) (*Record, error)
	Get(name string) (*Record, error)
	List() ([]*Record, error)
	Delete(name string) error
}
