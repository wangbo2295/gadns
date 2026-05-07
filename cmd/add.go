package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addIPs string

var addCmd = &cobra.Command{
	Use:   "add -i <ip_list> <name>",
	Short: "新增 DNS 记录",
	Long:  "为域名创建 A 记录，每个 IP 创建一条记录，多 IP 时自动分配权重。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if addIPs == "" {
			return fmt.Errorf("--ips flag is required")
		}

		cp, err := newProvider()
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}

		record, err := cp.Create(fullName(args[0]), strings.Split(addIPs, ","))
		if err != nil {
			return err
		}

		fmt.Printf("✓ add successful\n")
		fmt.Printf("  Name:  %s\n", record.Name)
		fmt.Printf("  CNAME: %s\n", record.CNAME)
		fmt.Printf("  IPs:   %s\n", strings.Join(record.IPs, ", "))
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addIPs, "ips", "i", "", "IP 地址列表（逗号分隔）")
	addCmd.MarkFlagRequired("ips")
	rootCmd.AddCommand(addCmd)
}
