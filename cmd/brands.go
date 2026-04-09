package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var brandsCmd = &cobra.Command{
	Use:   "brands",
	Short: "Manage brands",
}

var brandsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List brands",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewBrandService(client)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.BrandListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Slug", "Description"}
		rows := make([][]string, len(resp.Items))
		for i, b := range resp.Items {
			desc := b.Description
			if len(desc) > 50 {
				desc = desc[:50] + "..."
			}
			rows[i] = []string{
				fmt.Sprintf("%d", b.ID), b.Name, b.Slug, desc,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var brandsShowCmd = &cobra.Command{
	Use:   "show <id|slug>",
	Short: "Show brand details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewBrandService(client)
		brand, err := svc.GetByRef(cmdContext(), args[0])
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", brand.ID)},
			{Key: "Name", Value: brand.Name},
			{Key: "Slug", Value: brand.Slug},
			{Key: "Description", Value: brand.Description},
			{Key: "Voice", Value: brand.VoiceGuidelines},
			{Key: "Target Audience", Value: brand.TargetAudience},
			{Key: "Key Info", Value: brand.KeyInfo},
			{Key: "Contact", Value: brand.ContactInfo},
			{Key: "Mission", Value: brand.MissionStatement},
			{Key: "Values", Value: brand.Values},
			{Key: "Visual Identity", Value: brand.VisualIdentity},
			{Key: "Content Pillars", Value: brand.ContentPillars},
			{Key: "Competitors", Value: brand.Competitors},
			{Key: "Dos and Donts", Value: brand.DosAndDonts},
			{Key: "Created", Value: brand.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var brandsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a brand",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Brand name").Value(&name),
					huh.NewInput().Title("Slug (optional)").Value(&slug),
					huh.NewInput().Title("Description (optional)").Value(&description),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		client := mustClient()
		svc := api.NewBrandService(client)
		brand, err := svc.Create(cmdContext(), name, slug, description)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", brand.ID)},
			{Key: "Name", Value: brand.Name},
			{Key: "Slug", Value: brand.Slug},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var brandsUpdateCmd = &cobra.Command{
	Use:   "update <id|slug>",
	Short: "Update a brand",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flagMap := map[string]string{
			"name":            "name",
			"slug":            "slug",
			"description":     "description",
			"voice":           "voice_guidelines",
			"target-audience": "target_audience",
			"key-info":        "key_info",
			"contact":         "contact_info",
			"mission":         "mission_statement",
			"values":          "values",
			"visual-identity": "visual_identity",
			"content-pillars": "content_pillars",
			"competitors":     "competitors",
			"dos-and-donts":   "dos_and_donts",
		}

		fields := map[string]any{}
		for flag, jsonKey := range flagMap {
			if cmd.Flags().Changed(flag) {
				v, _ := cmd.Flags().GetString(flag)
				fields[jsonKey] = v
			}
		}

		if len(fields) == 0 {
			return fmt.Errorf("provide at least one flag to update")
		}

		client := mustClient()
		svc := api.NewBrandService(client)
		brand, err := svc.UpdateByRef(cmdContext(), args[0], fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", brand.ID)},
			{Key: "Name", Value: brand.Name},
			{Key: "Slug", Value: brand.Slug},
			{Key: "Description", Value: brand.Description},
			{Key: "Voice", Value: brand.VoiceGuidelines},
			{Key: "Target Audience", Value: brand.TargetAudience},
			{Key: "Key Info", Value: brand.KeyInfo},
			{Key: "Contact", Value: brand.ContactInfo},
			{Key: "Mission", Value: brand.MissionStatement},
			{Key: "Content Pillars", Value: brand.ContentPillars},
			{Key: "Created", Value: brand.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var brandsDeleteCmd = &cobra.Command{
	Use:   "delete <id|slug>",
	Short: "Delete a brand",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var confirm bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Delete brand %s?", args[0])).
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
		svc := api.NewBrandService(client)
		if err := svc.DeleteByRef(cmdContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Brand deleted.")
		return nil
	},
}

var brandsUseCmd = &cobra.Command{
	Use:   "use <id|slug>",
	Short: "Set default brand",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := mustLoadConfig()
		client := mustClient()
		svc := api.NewBrandService(client)
		brand, err := svc.GetByRef(cmdContext(), args[0])
		if err != nil {
			return err
		}

		c.DefaultBrandID = brand.ID
		c.DefaultBrandSlug = brand.Slug

		if err := c.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Default brand set to %s (%s)\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#00fff2")).Bold(true).Render(brand.Name),
			brand.Slug,
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(brandsCmd)

	brandsCmd.AddCommand(brandsListCmd)
	var page, perPage int
	addPaginationFlags(brandsListCmd, &page, &perPage)

	brandsCmd.AddCommand(brandsShowCmd)

	brandsCmd.AddCommand(brandsCreateCmd)
	brandsCreateCmd.Flags().String("name", "", "brand name")
	brandsCreateCmd.Flags().String("slug", "", "brand slug")
	brandsCreateCmd.Flags().String("description", "", "brand description")

	brandsCmd.AddCommand(brandsUpdateCmd)
	brandsUpdateCmd.Flags().String("name", "", "new name")
	brandsUpdateCmd.Flags().String("slug", "", "new slug")
	brandsUpdateCmd.Flags().String("description", "", "new description")
	brandsUpdateCmd.Flags().String("voice", "", "voice and tone guidelines")
	brandsUpdateCmd.Flags().String("target-audience", "", "target audience")
	brandsUpdateCmd.Flags().String("key-info", "", "key information")
	brandsUpdateCmd.Flags().String("contact", "", "point of contact")
	brandsUpdateCmd.Flags().String("mission", "", "mission statement")
	brandsUpdateCmd.Flags().String("values", "", "brand values")
	brandsUpdateCmd.Flags().String("visual-identity", "", "visual identity guidelines")
	brandsUpdateCmd.Flags().String("content-pillars", "", "content pillars")
	brandsUpdateCmd.Flags().String("competitors", "", "competitors")
	brandsUpdateCmd.Flags().String("dos-and-donts", "", "dos and donts")

	brandsCmd.AddCommand(brandsDeleteCmd)
	brandsCmd.AddCommand(brandsUseCmd)
}
