package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var ticketsCmd = &cobra.Command{
	Use:   "tickets",
	Short: "Manage support tickets",
}

var ticketsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tickets",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewTicketService(client)

		category, _ := cmd.Flags().GetString("category")
		priority, _ := cmd.Flags().GetString("priority")
		unresolved, _ := cmd.Flags().GetBool("unresolved")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.TicketListParams{
			Category: category, Priority: priority, Unresolved: unresolved,
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Subject", "Category", "Priority", "Status"}
		rows := make([][]string, len(resp.Items))
		for i, t := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", t.ID), t.Subject, t.Category, t.Priority, t.Status,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var ticketsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show ticket details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ticket ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewTicketService(client)
		ticket, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", ticket.ID)},
			{Key: "Subject", Value: ticket.Subject},
			{Key: "Description", Value: ticket.Description},
			{Key: "Category", Value: ticket.Category},
			{Key: "Priority", Value: ticket.Priority},
			{Key: "Status", Value: ticket.Status},
			{Key: "Created", Value: ticket.CreatedAt},
		}
		if ticket.ResolvedAt != "" {
			fields = append(fields, output.Field{Key: "Resolved At", Value: ticket.ResolvedAt})
			fields = append(fields, output.Field{Key: "Resolved Reason", Value: ticket.ResolvedReason})
		}
		if ticket.AdminNotes != "" {
			fields = append(fields, output.Field{Key: "Admin Notes", Value: ticket.AdminNotes})
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var ticketsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a ticket",
	RunE: func(cmd *cobra.Command, args []string) error {
		subject, _ := cmd.Flags().GetString("subject")
		description, _ := cmd.Flags().GetString("description")
		category, _ := cmd.Flags().GetString("category")
		priority, _ := cmd.Flags().GetString("priority")

		if subject == "" {
			if !isInteractiveTerminal() {
				return fmt.Errorf("--subject is required")
			}
			if category == "" {
				category = "bug_report"
			}
			if priority == "" {
				priority = "medium"
			}
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Subject").Value(&subject),
					huh.NewText().Title("Description").Value(&description),
					huh.NewSelect[string]().
						Title("Category").
						Options(
							huh.NewOption("Bug Report", "bug_report"),
							huh.NewOption("Feature Request", "feature_request"),
						).
						Value(&category),
					huh.NewSelect[string]().
						Title("Priority").
						Options(
							huh.NewOption("Low", "low"),
							huh.NewOption("Medium", "medium"),
							huh.NewOption("High", "high"),
						).
						Value(&priority),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		client := mustClient()
		svc := api.NewTicketService(client)
		ticket, err := svc.Create(cmdContext(), subject, description, category, priority)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", ticket.ID)},
			{Key: "Subject", Value: ticket.Subject},
			{Key: "Status", Value: ticket.Status},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var ticketsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a ticket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ticket ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("subject") {
			v, _ := cmd.Flags().GetString("subject")
			fields["subject"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			fields["description"] = v
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetString("priority")
			fields["priority"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --subject, --description, --priority)")
		}

		client := mustClient()
		svc := api.NewTicketService(client)
		ticket, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", ticket.ID)},
			{Key: "Subject", Value: ticket.Subject},
			{Key: "Priority", Value: ticket.Priority},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ticketsCmd)

	ticketsCmd.AddCommand(ticketsListCmd)
	ticketsListCmd.Flags().String("category", "", "filter by category (bug_report, feature_request)")
	ticketsListCmd.Flags().String("priority", "", "filter by priority (low, medium, high)")
	ticketsListCmd.Flags().Bool("unresolved", false, "show only unresolved tickets")
	var page, perPage int
	addPaginationFlags(ticketsListCmd, &page, &perPage)

	ticketsCmd.AddCommand(ticketsShowCmd)

	ticketsCmd.AddCommand(ticketsCreateCmd)
	ticketsCreateCmd.Flags().String("subject", "", "ticket subject")
	ticketsCreateCmd.Flags().String("description", "", "ticket description")
	ticketsCreateCmd.Flags().String("category", "", "category: bug_report, feature_request")
	ticketsCreateCmd.Flags().String("priority", "", "priority: low, medium, high")

	ticketsCmd.AddCommand(ticketsUpdateCmd)
	ticketsUpdateCmd.Flags().String("subject", "", "new subject")
	ticketsUpdateCmd.Flags().String("description", "", "new description")
	ticketsUpdateCmd.Flags().String("priority", "", "new priority")
}
