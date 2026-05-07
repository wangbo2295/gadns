package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 DNS 记录",
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := newProvider()
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}

		records, err := cp.List()
		if err != nil {
			return err
		}

		if len(records) == 0 {
			fmt.Println("No records found")
			return nil
		}

		fmt.Printf("Total Records: %d\n\n", len(records))
		fmt.Printf("%-30s | %-40s | %s\n", "Name", "CNAME", "IPs")
		fmt.Printf("%s-+-%s-+-%s\n",
			strings.Repeat("-", 30), strings.Repeat("-", 40), strings.Repeat("-", 40))

		for _, r := range records {
			ips := strings.Join(r.IPs, ", ")
			if len(ips) > 37 {
				ips = ips[:37] + "..."
			}
			fmt.Printf("%-30s | %-40s | %s\n", r.Name, r.CNAME, ips)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
