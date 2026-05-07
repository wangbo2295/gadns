package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除 DNS 记录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := newProvider()
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}

		if err := cp.Delete(fullName(args[0])); err != nil {
			return err
		}

		fmt.Printf("✓ Record '%s' deleted\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
