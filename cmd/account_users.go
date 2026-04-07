package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage team members",
}

var membersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List team members",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAccountService(client)
		members, err := svc.ListMembers(cmdContext())
		if err != nil {
			return err
		}

		headers := []string{"ID", "Email", "Roles", "Created"}
		rows := make([][]string, len(members))
		for i, m := range members {
			rows[i] = []string{
				fmt.Sprintf("%d", m.ID), m.Email,
				strings.Join(m.Roles, ", "), m.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, nil))
		return nil
	},
}

var membersShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show member details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid member ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		member, brandIDs, err := svc.GetMember(cmdContext(), id)
		if err != nil {
			return err
		}

		brandIDStrs := make([]string, len(brandIDs))
		for i, bid := range brandIDs {
			brandIDStrs[i] = fmt.Sprintf("%d", bid)
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", member.ID)},
			{Key: "Email", Value: member.Email},
			{Key: "Roles", Value: strings.Join(member.Roles, ", ")},
			{Key: "Brand Access", Value: strings.Join(brandIDStrs, ", ")},
			{Key: "Created", Value: member.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var membersUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update member roles",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid member ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("admin") {
			v, _ := cmd.Flags().GetBool("admin")
			fields["admin"] = v
		}
		if cmd.Flags().Changed("editor") {
			v, _ := cmd.Flags().GetBool("editor")
			fields["editor"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no roles to update (use --admin, --editor)")
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		member, err := svc.UpdateMember(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", member.ID)},
			{Key: "Email", Value: member.Email},
			{Key: "Roles", Value: strings.Join(member.Roles, ", ")},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var membersDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Remove a team member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid member ID: %s", args[0])
		}

		var confirm bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Remove member %d?", id)).
			Description("They will lose access to this account.").
			Value(&confirm).
			Run(); err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		if err := svc.DeleteMember(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Member removed.")
		return nil
	},
}

var membersBrandAccessCmd = &cobra.Command{
	Use:   "brand-access <id>",
	Short: "Update brand access for a member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid member ID: %s", args[0])
		}

		brandIDsStr, _ := cmd.Flags().GetString("brand-ids")
		if brandIDsStr == "" {
			return fmt.Errorf("--brand-ids is required (comma-separated)")
		}

		brandIDs, err := parseIntList(brandIDsStr)
		if err != nil {
			return fmt.Errorf("invalid brand IDs: %w", err)
		}

		client := mustClient()
		svc := api.NewAccountService(client)
		member, err := svc.UpdateBrandAccess(cmdContext(), id, brandIDs)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", member.ID)},
			{Key: "Email", Value: member.Email},
		}
		fmt.Print(formatter.FormatItem(out))
		fmt.Println("Brand access updated.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(membersCmd)

	membersCmd.AddCommand(membersListCmd)
	membersCmd.AddCommand(membersShowCmd)

	membersCmd.AddCommand(membersUpdateCmd)
	membersUpdateCmd.Flags().Bool("admin", false, "admin role")
	membersUpdateCmd.Flags().Bool("editor", false, "editor role")

	membersCmd.AddCommand(membersDeleteCmd)

	membersCmd.AddCommand(membersBrandAccessCmd)
	membersBrandAccessCmd.Flags().String("brand-ids", "", "comma-separated brand IDs")
}
