package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var updateIPs string

var updateCmd = &cobra.Command{
	Use:   "update -ips <ip_list> <name>",
	Short: "更新 DNS 记录",
	Long:  "更新域名对应的 IP 映射，先删除旧记录再重新创建",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateIPs == "" {
			return fmt.Errorf("--ips flag is required")
		}

		cp, err := newProvider()
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}

		record, err := cp.Update(fullName(args[0]), strings.Split(updateIPs, ","))
		if err != nil {
			return err
		}

		fmt.Printf("✓ update successful\n")
		fmt.Printf("  Name:  %s\n", record.Name)
		fmt.Printf("  CNAME: %s\n", record.CNAME)
		fmt.Printf("  IPs:   %s\n", strings.Join(record.IPs, ", "))
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVarP(&updateIPs, "ips", "i", "", "IP 地址列表（逗号分隔）")
	updateCmd.MarkFlagRequired("ips")
	rootCmd.AddCommand(updateCmd)
}
