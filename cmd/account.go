package cmd

import (
	"fmt"
	"strings"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage account settings",
}

var accountShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show account details",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAccountService(client)
		detail, err := svc.Get(cmdContext())
		if err != nil {
			return err
		}

		a := detail.Account
		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", a.ID)},
			{Key: "Name", Value: a.Name},
			{Key: "Slug", Value: a.Slug},
			{Key: "Members", Value: fmt.Sprintf("%d", len(detail.Users))},
			{Key: "Created", Value: a.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(detail.Users) > 0 {
			fmt.Println()
			headers := []string{"ID", "Email", "Roles"}
			rows := make([][]string, len(detail.Users))
			for i, u := range detail.Users {
				rows[i] = []string{
					fmt.Sprintf("%d", u.ID), u.Email, strings.Join(u.Roles, ", "),
				}
			}
			fmt.Print(formatter.FormatList(headers, rows, nil))
		}

		return nil
	},
}

var accountUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update account settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")

		if name == "" && slug == "" {
			return fmt.Errorf("provide --name or --slug to update")
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		account, err := svc.Update(cmdContext(), name, slug)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", account.ID)},
			{Key: "Name", Value: account.Name},
			{Key: "Slug", Value: account.Slug},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.AddCommand(accountShowCmd)

	accountCmd.AddCommand(accountUpdateCmd)
	accountUpdateCmd.Flags().String("name", "", "account name")
	accountUpdateCmd.Flags().String("slug", "", "account slug")
}
