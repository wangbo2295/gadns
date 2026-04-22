// types/record.go
package types

// Record DNS记录信息
type Record struct {
	Name  string
	CNAME string
	IPs   []string
}
