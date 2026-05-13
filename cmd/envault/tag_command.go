package main

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/envault/internal/vault"
	"github.com/spf13/cobra"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on env keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <key> <tag1,tag2,...>",
		Short: "Set tags on a key (comma-separated)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}
			tags := strings.Split(args[1], ",")
			if err := vault.SetTags(v, args[0], tags); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tags set on %s: %s\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get tags for a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}
			tags, err := vault.GetTags(v, args[0])
			if err != nil {
				return err
			}
			if len(tags) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%s has no tags\n", args[0])
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", args[0], strings.Join(tags, ", "))
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <tag>",
		Short: "List all keys with a given tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := vault.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			v, err := vault.New(cfg)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}
			results, err := vault.ListByTag(v, args[0])
			if err != nil {
				return err
			}
			if len(results) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no keys tagged with %q\n", args[0])
				return nil
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]\n", r.Key, strings.Join(r.Tags, ", "))
			}
			return nil
		},
	}

	tagCmd.AddCommand(setCmd, getCmd, listCmd)
	rootCmd.AddCommand(tagCmd)
}
