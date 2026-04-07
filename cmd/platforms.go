package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var platformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "Manage platforms",
}

var platformsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewPlatformService(client)

		active, _ := cmd.Flags().GetString("active")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.PlatformListParams{
			Active: active, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Type", "Slug", "Active", "Model"}
		rows := make([][]string, len(resp.Items))
		for i, p := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", p.ID), p.Name, p.Type, p.Slug,
				fmt.Sprintf("%t", p.Active), p.ModelDisplayName(),
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var platformsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show platform details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		platform, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", platform.ID)},
			{Key: "Name", Value: platform.Name},
			{Key: "Type", Value: platform.Type},
			{Key: "Slug", Value: platform.Slug},
			{Key: "Active", Value: fmt.Sprintf("%t", platform.Active)},
			{Key: "Model", Value: platform.ModelDisplayName()},
			{Key: "Created", Value: platform.CreatedAt},
		}
		if platform.PromptTemplate != "" {
			v := platform.PromptTemplate
			if outFormat == "table" && len(v) > 80 {
				v = v[:80] + "..."
			}
			fields = append(fields, output.Field{Key: "Prompt Template", Value: v})
		}
		fmt.Print(formatter.FormatItem(fields))
		if outFormat == "table" && len(platform.PromptTemplate) > 80 {
			fmt.Printf("\nPrompt Template:\n%s\n", platform.PromptTemplate)
		}
		return nil
	},
}

var platformsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		platformType, _ := cmd.Flags().GetString("type")
		slug, _ := cmd.Flags().GetString("slug")

		if name == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Platform name").Value(&name),
					huh.NewInput().Title("Type (e.g. Platforms::Twitter)").Value(&platformType),
					huh.NewInput().Title("Slug (optional)").Value(&slug),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		model, _ := cmd.Flags().GetString("model")

		client := mustClient()
		svc := api.NewPlatformService(client)
		platform, err := svc.Create(cmdContext(), name, platformType, slug, true, "", model)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", platform.ID)},
			{Key: "Name", Value: platform.Name},
			{Key: "Type", Value: platform.Type},
			{Key: "Slug", Value: platform.Slug},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var platformsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a platform",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			fields["name"] = v
		}
		if cmd.Flags().Changed("active") {
			v, _ := cmd.Flags().GetBool("active")
			fields["active"] = v
		}
		if cmd.Flags().Changed("slug") {
			v, _ := cmd.Flags().GetString("slug")
			fields["slug"] = v
		}
		if cmd.Flags().Changed("model") {
			v, _ := cmd.Flags().GetString("model")
			fields["model_id"] = v
		}
		if cmd.Flags().Changed("prompt") {
			v, _ := cmd.Flags().GetString("prompt")
			fields["prompt_template"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --name, --active, --slug, --model, --prompt)")
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		platform, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", platform.ID)},
			{Key: "Name", Value: platform.Name},
			{Key: "Active", Value: fmt.Sprintf("%t", platform.Active)},
			{Key: "Model", Value: platform.ModelDisplayName()},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var platformsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a platform",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}

		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Delete platform %d?", id)).
			Description("This cannot be undone.").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Platform deleted.")
		return nil
	},
}

// Label subcommands

var platformsLabelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage platform labels",
}

var platformsLabelsListCmd = &cobra.Command{
	Use:   "list <platform_id>",
	Short: "List labels for a platform",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		client := mustClient()
		svc := api.NewPlatformService(client)
		resp, err := svc.ListLabels(cmdContext(), platformID, api.PlatformListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Slug", "Description"}
		rows := make([][]string, len(resp.Items))
		for i, l := range resp.Items {
			desc := l.Description
			if len(desc) > 50 {
				desc = desc[:50] + "..."
			}
			rows[i] = []string{
				fmt.Sprintf("%d", l.ID), l.Name, l.Slug, desc,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var platformsLabelsShowCmd = &cobra.Command{
	Use:   "show <platform_id> <label_id>",
	Short: "Show label details",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}
		labelID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid label ID: %s", args[1])
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		label, err := svc.GetLabel(cmdContext(), platformID, labelID)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", label.ID)},
			{Key: "Name", Value: label.Name},
			{Key: "Slug", Value: label.Slug},
			{Key: "Description", Value: label.Description},
			{Key: "Created", Value: label.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var platformsLabelsCreateCmd = &cobra.Command{
	Use:   "create <platform_id>",
	Short: "Create a label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		fields := map[string]any{"name": name}
		if cmd.Flags().Changed("slug") {
			v, _ := cmd.Flags().GetString("slug")
			fields["slug"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			fields["description"] = v
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		label, err := svc.CreateLabel(cmdContext(), platformID, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", label.ID)},
			{Key: "Name", Value: label.Name},
			{Key: "Slug", Value: label.Slug},
			{Key: "Description", Value: label.Description},
			{Key: "Created", Value: label.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var platformsLabelsUpdateCmd = &cobra.Command{
	Use:   "update <platform_id> <label_id>",
	Short: "Update a label",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}
		labelID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid label ID: %s", args[1])
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
			return fmt.Errorf("provide at least one flag to update")
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		label, err := svc.UpdateLabel(cmdContext(), platformID, labelID, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", label.ID)},
			{Key: "Name", Value: label.Name},
			{Key: "Slug", Value: label.Slug},
			{Key: "Description", Value: label.Description},
			{Key: "Created", Value: label.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var platformsLabelsDeleteCmd = &cobra.Command{
	Use:   "delete <platform_id> <label_id>",
	Short: "Delete a label",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid platform ID: %s", args[0])
		}
		labelID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid label ID: %s", args[1])
		}

		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Delete label %d?", labelID)).
			Description("This cannot be undone.").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewPlatformService(client)
		if err := svc.DeleteLabel(cmdContext(), platformID, labelID); err != nil {
			return err
		}
		fmt.Println("Label deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(platformsCmd)

	platformsCmd.AddCommand(platformsListCmd)
	platformsListCmd.Flags().String("active", "", "filter by active status (true/false)")
	var page, perPage int
	addPaginationFlags(platformsListCmd, &page, &perPage)

	platformsCmd.AddCommand(platformsShowCmd)

	platformsCmd.AddCommand(platformsCreateCmd)
	platformsCreateCmd.Flags().String("name", "", "platform name")
	platformsCreateCmd.Flags().String("type", "", "platform type (e.g. Platforms::Twitter)")
	platformsCreateCmd.Flags().String("slug", "", "platform slug")
	platformsCreateCmd.Flags().String("model", "", "AI model provider ID (e.g. claude-sonnet-4-5-20241022)")

	platformsCmd.AddCommand(platformsUpdateCmd)
	platformsUpdateCmd.Flags().String("name", "", "new name")
	platformsUpdateCmd.Flags().Bool("active", true, "active status")
	platformsUpdateCmd.Flags().String("slug", "", "new slug")
	platformsUpdateCmd.Flags().String("model", "", "AI model provider ID (e.g. claude-sonnet-4-5-20241022)")
	platformsUpdateCmd.Flags().String("prompt", "", "new prompt template")

	platformsCmd.AddCommand(platformsDeleteCmd)

	platformsCmd.AddCommand(platformsLabelsCmd)
	platformsLabelsCmd.AddCommand(platformsLabelsListCmd)
	var labelPage, labelPerPage int
	addPaginationFlags(platformsLabelsListCmd, &labelPage, &labelPerPage)
	platformsLabelsCmd.AddCommand(platformsLabelsShowCmd)
	platformsLabelsCmd.AddCommand(platformsLabelsCreateCmd)
	platformsLabelsCreateCmd.Flags().String("name", "", "label name")
	platformsLabelsCreateCmd.Flags().String("slug", "", "label slug")
	platformsLabelsCreateCmd.Flags().String("description", "", "label description")
	platformsLabelsCmd.AddCommand(platformsLabelsUpdateCmd)
	platformsLabelsUpdateCmd.Flags().String("name", "", "new name")
	platformsLabelsUpdateCmd.Flags().String("slug", "", "new slug")
	platformsLabelsUpdateCmd.Flags().String("description", "", "new description")
	platformsLabelsCmd.AddCommand(platformsLabelsDeleteCmd)
}
