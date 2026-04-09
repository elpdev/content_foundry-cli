package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage campaigns",
}

var campaignsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List campaigns",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewCampaignService(client)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.CampaignListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Slug", "Description"}
		rows := make([][]string, len(resp.Items))
		for i, campaign := range resp.Items {
			description := campaign.Description
			if len(description) > 50 {
				description = description[:50] + "..."
			}
			rows[i] = []string{
				fmt.Sprintf("%d", campaign.ID), campaign.Name, campaign.Slug, description,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var campaignsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show campaign details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid campaign ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewCampaignService(client)
		if outFormat == "json" {
			raw, err := svc.GetRaw(cmdContext(), id)
			if err != nil {
				return err
			}
			fmt.Print(formatter.FormatRaw(raw))
			return nil
		}

		campaign, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", campaign.ID)},
			{Key: "Name", Value: campaign.Name},
			{Key: "Slug", Value: campaign.Slug},
			{Key: "Description", Value: campaign.Description},
			{Key: "Created", Value: campaign.CreatedAt},
			{Key: "Updated", Value: campaign.UpdatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var campaignsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		description, _ := cmd.Flags().GetString("description")

		if name == "" || slug == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Campaign name").Value(&name),
					huh.NewInput().Title("Slug").Value(&slug),
					huh.NewInput().Title("Description (optional)").Value(&description),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		client := mustClient()
		svc := api.NewCampaignService(client)
		campaign, err := svc.Create(cmdContext(), name, slug, description)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", campaign.ID)},
			{Key: "Name", Value: campaign.Name},
			{Key: "Slug", Value: campaign.Slug},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var campaignsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid campaign ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			fields["name"] = v
		}
		if cmd.Flags().Changed("slug") {
			v, _ := cmd.Flags().GetString("slug")
			fields["slug"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			fields["description"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --name, --slug, --description)")
		}

		client := mustClient()
		svc := api.NewCampaignService(client)
		campaign, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", campaign.ID)},
			{Key: "Name", Value: campaign.Name},
			{Key: "Slug", Value: campaign.Slug},
			{Key: "Description", Value: campaign.Description},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var campaignsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid campaign ID: %s", args[0])
		}

		var confirm bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Delete campaign %d?", id)).
			Description("This cannot be undone.").
			Value(&confirm).
			Run(); err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewCampaignService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Campaign deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(campaignsCmd)

	campaignsCmd.AddCommand(campaignsListCmd)
	var page, perPage int
	addPaginationFlags(campaignsListCmd, &page, &perPage)

	campaignsCmd.AddCommand(campaignsShowCmd)

	campaignsCmd.AddCommand(campaignsCreateCmd)
	campaignsCreateCmd.Flags().String("name", "", "campaign name")
	campaignsCreateCmd.Flags().String("slug", "", "campaign slug")
	campaignsCreateCmd.Flags().String("description", "", "campaign description")

	campaignsCmd.AddCommand(campaignsUpdateCmd)
	campaignsUpdateCmd.Flags().String("name", "", "new campaign name")
	campaignsUpdateCmd.Flags().String("slug", "", "new campaign slug")
	campaignsUpdateCmd.Flags().String("description", "", "new campaign description")

	campaignsCmd.AddCommand(campaignsDeleteCmd)
}
