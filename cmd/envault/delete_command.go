package main

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envault/internal/vault"
	"github.com/spf13/cobra"
)

func init() {
	var mustExist bool

	cmd := &cobra.Command{
		Use:   "delete <KEY> [KEY...]",
		Short: "Delete one or more keys from the env file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := vault.LoadConfig(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("init vault: %w", err)
			}

			res, err := v.DeleteKeys(vault.DeleteOptions{
				Keys:      args,
				MustExist: mustExist,
			})
			if err != nil {
				return err
			}

			for _, k := range res.Deleted {
				fmt.Fprintf(os.Stdout, "deleted: %s\n", k)
			}
			for _, k := range res.NotFound {
				fmt.Fprintf(os.Stdout, "not found (skipped): %s\n", k)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&mustExist, "must-exist", false, "fail if any key does not exist")
	rootCmd.AddCommand(cmd)
}
