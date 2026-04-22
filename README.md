# SmartDNS

SmartDNS 是一个 Go 语言的 DNS 管理组件库和命令行工具，根据 IP 集合生成 CNAME 记录，支持腾讯云智能DNS和本地调试两种实现。

## 特性

- 支持多种IP格式：单IP、IP范围、CIDR段
- 接口与实现解耦，易于扩展
- 内置腾讯云智能DNS和本地调试两种实现
- **腾讯云智能调度**：支持按地域、运营商智能返回不同IP
- 提供完整的命令行工具
- 严格格式校验
- 完整的测试覆盖

## 安装

### 作为库使用

```bash
go get github.com/yourusername/smartdns
```

### 安装命令行工具

```bash
# 从源码构建
go build -o smartdns cmd/smartdns/main.go
sudo mv smartdns /usr/local/bin/

# 或使用 go install
go install github.com/yourusername/smartdns/cmd/smartdns@latest
```

## 命令行工具使用

### 配置文件

**本地模式配置** (`~/.smartdns/local.yaml`):

```yaml
hosts_path: "/etc/hosts"
storage_path: "~/.smartdns/records.json"
domain: "smartdns.local"
```

**腾讯云模式配置** (`~/.smartdns/tencent.yaml`):

```yaml
secret_id: "your_secret_id"
secret_key: "your_secret_key"
region: "ap-guangzhou"
domain: "example.com"
```

### 命令示例

```bash
# 新增记录
smartdns add -ips 1.1.1.1 app

# 多个 IP
smartdns add -ips 1.1.1.1,2.2.2.2 web

# IP 范围
smartdns add -ips 10.0.0.1-10.0.0.10 api

# CIDR 段
smartdns add -ips 192.168.1.0/24 db

# 混合格式
smartdns add -ips "1.1.1.1,10.0.0.1-10.0.0.5,192.168.1.0/30" cache

# 查询记录
smartdns get app

# 列出所有记录
smartdns list

# 更新记录
smartdns update -ips 5.5.5.5 app

# 删除记录
smartdns delete app

# 使用自定义配置
smartdns add -ips 1.1.1.1 -c /path/to/config.yaml app
```

### 选项

- `-ips <list>`: IP 地址列表（逗号分隔）[add/update 必需]
- `-c <path>`: 配置文件路径 (默认: ~/.smartdns/local.yaml)
- `-h, --help`: 显示帮助信息

### 命令帮助

```bash
smartdns help
```

## 作为库使用

### 使用工厂函数（推荐）

```go
package main

import (
    "fmt"
    "github.com/yourusername/smartdns/provider"
)

func main() {
    dns, err := provider.New("local", "~/.smartdns/local.yaml")
    if err != nil {
        panic(err)
    }

    record, err := dns.Create("app", []string{
        "1.1.1.1",         // 单IP
        "1.1.1.5-1.1.1.7", // IP范围
        "2.2.2.0/30",      // CIDR段
    })

    fmt.Printf("CNAME: %s\n", record.CNAME)
}
```

### 直接使用 local provider

```go
package main

import (
    "fmt"
    "github.com/yourusername/smartdns/provider/local"
)

func main() {
    provider := local.NewProvider(&local.Config{
        HostsPath:   "/etc/hosts",
        StoragePath: "~/.smartdns/records.json",
        Domain:      "smartdns.local",
    })

    record, err := provider.Create("web", []string{"10.0.0.1-10.0.0.5"})
    if err != nil {
        panic(err)
    }

    fmt.Printf("CNAME: %s\n", record.CNAME)
}
```

## IP格式支持

| 格式 | 示例 | 说明 |
|------|------|------|
| 单IP | `1.1.1.1` | 单个 IPv4 地址 |
| IP范围 | `1.1.1.1-1.1.1.10` | 起始到结束的 IP 范围 |
| CIDR段 | `1.1.1.0/24` | CIDR 表示法 |
| 混合 | `["1.1.1.1", "1.1.1.5-1.1.1.10", "2.2.2.0/30"]` | 多种格式组合 |

## 接口

```go
type SmartDNS interface {
    Create(name string, ips []string) (*Record, error)
    Update(name string, ips []string) (*Record, error)
    Get(name string) (*Record, error)
    List() ([]*Record, error)
    Delete(name string) error
}

type Record struct {
    Name  string   // 记录名称
    CNAME string   // 生成的CNAME
    IPs   []string // 原始IP列表（支持列表、范围、段格式）
}
```

## 项目结构

```
smartdns/
├── main.go              # CLI 工具入口 (package main)
├── core/                # 核心接口和类型 (package core)
│   ├── smartdns.go      # SmartDNS 接口和 Record 类型
│   └── smartdns_test.go
├── provider/            # Provider 实现 (package provider)
│   ├── factory.go       # 工厂函数
│   ├── local/           # 本地实现
│   │   ├── provider.go
│   │   ├── hosts.go
│   │   └── storage.go
│   └── tencent/         # 腾讯云实现
│       ├── provider.go
│       └── client.go
├── iputil/              # IP 处理工具 (package iputil)
│   ├── parser.go        # IP 解析
│   └── validator.go     # 格式校验
├── config/              # 配置加载 (package config)
│   └── loader.go        # YAML 配置加载
├── Makefile             # 构建脚本
└── README.md
```

## 腾讯云智能DNS

腾讯云模式支持**智能调度**功能，可以为不同的用户群体返回不同的IP地址。

### 智能调度特性

- **地域调度**：支持按省份/地区返回不同IP（如：北京用户返回IP1，上海用户返回IP2）
- **运营商调度**：支持按运营商返回不同IP（如：电信用户返回IP1，联通用户返回IP2）
- **默认线路**：为未匹配地域或运营商的用户提供默认IP

### 智能调度配置

当前内置默认配置：

```go
DefaultSmartRoutingConfig = SmartRoutingConfig{
    Enabled: true,
    Regions: []string{
        "北京", "上海", "广东", "江苏", "浙江", "四川", "湖北", "福建",
    },
    Carriers: []string{
        "电信", "联通", "移动",
    },
}
```

### 工作原理

当启用智能调度时，系统会为每个IP创建多条DNS记录：

1. **默认线路记录**：第一条IP作为默认返回值
2. **地域线路记录**：为配置的每个地区创建一条记录
3. **运营商线路记录**：为每个运营商创建一条记录

例如，使用单个IP `1.1.1.1` 创建记录 `app`：
- 创建 `app.example.com` → `1.1.1.1`（默认线路）
- 创建 `app.example.com` → `1.1.1.1`（北京线路）
- 创建 `app.example.com` → `1.1.1.1`（上海线路）
- ...
- 创建 `app.example.com` → `1.1.1.1`（电信线路）
- 创建 `app.example.com` → `1.1.1.1`（联通线路）
- ...

### 使用示例

```bash
# 创建智能调度记录（自动为每个IP创建多条线路记录）
smartdns -c ~/.smartdns/tencent.yaml add -ips 1.1.1.1 app

# 使用多个IP（每个IP都会创建地域和运营商线路记录）
smartdns -c ~/.smartdns/tencent.yaml add -ips 1.1.1.1,2.2.2.2 api

# 更新记录（删除所有旧记录，重新创建）
smartdns -c ~/.smartdns/tencent.yaml update -ips 3.3.3.3 app

# 查询记录（聚合显示所有IP）
smartdns -c ~/.smartdns/tencent.yaml get app

# 列出所有记录
smartdns -c ~/.smartdns/tencent.yaml list

# 删除记录（删除所有相关线路记录）
smartdns -c ~/.smartdns/tencent.yaml delete app
```

### 腾讯云配置示例

创建 `~/.smartdns/tencent.yaml`:

```yaml
secret_id: "AKIDxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
secret_key: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
region: "ap-guangzhou"  # 可选，默认为 ap-guangzhou
domain: "example.com"    # 你的域名
```

### 注意事项

1. **API密钥安全**：请妥善保管腾讯云API密钥，建议使用子账号并限制权限
2. **记录数量**：智能调度会创建多条记录，注意DNSPod套餐的记录数量限制
3. **生效时间**：DNS记录修改后通常需要几分钟生效
4. **TTL设置**：默认TTL为600秒（10分钟）

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./provider/local
go test ./iputil

# 查看测试覆盖率
go test -cover ./...

# 运行 CLI 测试
go test ./cmd
```

## 开发

```bash
# 安装依赖
go mod tidy

# 格式化代码
go fmt ./...

# 静态检查
go vet ./...
```

## License

MIT
