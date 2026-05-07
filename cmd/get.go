package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <domain>",
	Short: "查询 DNS 记录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := newProvider()
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}

		record, err := cp.Get(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Name:   %s\n", record.Name)
		fmt.Printf("CNAME:  %s\n", record.CNAME)
		fmt.Printf("IPs:    %s\n", strings.Join(record.IPs, ", "))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
