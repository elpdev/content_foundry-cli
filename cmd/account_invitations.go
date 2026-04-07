package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var invitationsCmd = &cobra.Command{
	Use:   "invitations",
	Short: "Manage account invitations",
}

var invitationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pending invitations",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAccountService(client)
		invitations, err := svc.ListInvitations(cmdContext())
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Email", "Roles", "Created"}
		rows := make([][]string, len(invitations))
		for i, inv := range invitations {
			rows[i] = []string{
				fmt.Sprintf("%d", inv.ID), inv.Name, inv.Email,
				strings.Join(inv.Roles, ", "), inv.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, nil))
		return nil
	},
}

var invitationsShowCmd = &cobra.Command{
	Use:   "show <token>",
	Short: "Show invitation details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAccountService(client)
		inv, err := svc.GetInvitation(cmdContext(), args[0])
		if err != nil {
			return err
		}

		brandIDStrs := make([]string, len(inv.BrandIDs))
		for i, bid := range inv.BrandIDs {
			brandIDStrs[i] = fmt.Sprintf("%d", bid)
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", inv.ID)},
			{Key: "Name", Value: inv.Name},
			{Key: "Email", Value: inv.Email},
			{Key: "Token", Value: inv.Token},
			{Key: "Roles", Value: strings.Join(inv.Roles, ", ")},
			{Key: "Brand IDs", Value: strings.Join(brandIDStrs, ", ")},
			{Key: "Created", Value: inv.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var invitationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invitation",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")

		if name == "" || email == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Name").Value(&name),
					huh.NewInput().Title("Email").Value(&email),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		fields := map[string]any{
			"name":  name,
			"email": email,
		}
		if cmd.Flags().Changed("admin") {
			v, _ := cmd.Flags().GetBool("admin")
			fields["admin"] = v
		}
		if cmd.Flags().Changed("brand-ids") {
			brandIDsStr, _ := cmd.Flags().GetString("brand-ids")
			brandIDs, err := parseIntList(brandIDsStr)
			if err != nil {
				return fmt.Errorf("invalid brand IDs: %w", err)
			}
			fields["brand_ids"] = brandIDs
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		inv, err := svc.CreateInvitation(cmdContext(), fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", inv.ID)},
			{Key: "Name", Value: inv.Name},
			{Key: "Email", Value: inv.Email},
			{Key: "Token", Value: inv.Token},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var invitationsDeleteCmd = &cobra.Command{
	Use:   "delete <token>",
	Short: "Delete an invitation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var confirm bool
		huh.NewConfirm().
			Title("Delete this invitation?").
			Description("The invite link will no longer work.").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		if err := svc.DeleteInvitation(cmdContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Invitation deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(invitationsCmd)

	invitationsCmd.AddCommand(invitationsListCmd)
	invitationsCmd.AddCommand(invitationsShowCmd)

	invitationsCmd.AddCommand(invitationsCreateCmd)
	invitationsCreateCmd.Flags().String("name", "", "invitee name")
	invitationsCreateCmd.Flags().String("email", "", "invitee email")
	invitationsCreateCmd.Flags().Bool("admin", false, "grant admin role")
	invitationsCreateCmd.Flags().String("brand-ids", "", "comma-separated brand IDs to grant access to")

	invitationsCmd.AddCommand(invitationsDeleteCmd)
}
