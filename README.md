# SmartDNS

SmartDNS 是一个 Go 语言的通用 DNS 管理组件库，根据 IP 集合生成 CNAME 记录。

## 特性

- 支持多种IP格式：单IP、IP范围、CIDR段
- 接口与实现解耦，易于扩展
- 内置腾讯云智能DNS和本地调试两种实现
- 严格格式校验
- 完整的测试覆盖

## 安装

```bash
go get github.com/yourusername/smartdns
```

## 使用

### 本地实现（开发调试）

创建配置文件 `~/.smartdns/local.yaml`:

```yaml
hosts_path: "/etc/hosts"
storage_path: "~/.smartdns/records.json"
domain: "smartdns.local"
```

```go
package main

import (
    "fmt"
    "github.com/yourusername/smartdns/smartdns"
)

func main() {
    dns, _ := smartdns.New("local", "~/.smartdns/local.yaml")

    // 创建记录
    record, _ := dns.Create("app", []string{
        "1.1.1.1",
        "1.1.1.5-1.1.1.10",
        "2.2.2.0/30",
    })

    fmt.Println(record.CNAME) // app.smartdns.local
}
```

### 腾讯云实现

创建配置文件 `~/.smartdns/tencent.yaml`:

```yaml
secret_id: "your_secret_id"
secret_key: "your_secret_key"
region: "ap-guangzhou"
domain: "example.com"
```

```go
dns, _ := smartdns.New("tencent", "~/.smartdns/tencent.yaml")
record, _ := dns.Create("app", []string{"1.1.1.1", "2.2.2.2"})
```

## IP格式支持

| 格式 | 示例 |
|------|------|
| 单IP | `1.1.1.1` |
| IP范围 | `1.1.1.5-1.1.1.10` |
| CIDR段 | `1.1.1.0/24` |

## 接口

```go
type SmartDNS interface {
    Create(name string, ips []string) (*Record, error)
    Update(name string, ips []string) (*Record, error)
    Get(name string) (*Record, error)
    List() ([]*Record, error)
    Delete(name string) error
}
```

## 运行测试

```bash
go test ./...
```

## License

MIT
