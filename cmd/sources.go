package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Manage content sources",
}

var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewSourceService(client)

		active, _ := cmd.Flags().GetString("active")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.SourceListParams{
			Active: active, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Type", "Active", "Last Fetched"}
		rows := make([][]string, len(resp.Items))
		for i, s := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", s.ID), s.Name, s.Type,
				fmt.Sprintf("%t", s.Active), s.LastFetchedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var sourcesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show source details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewSourceService(client)
		if outFormat == "json" {
			raw, err := svc.GetRaw(cmdContext(), id)
			if err != nil {
				return err
			}
			fmt.Print(formatter.FormatRaw(raw))
			return nil
		}

		source, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		cfgJSON, _ := json.Marshal(source.Config)

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", source.ID)},
			{Key: "Name", Value: source.Name},
			{Key: "Type", Value: source.Type},
			{Key: "Active", Value: fmt.Sprintf("%t", source.Active)},
			{Key: "Polling", Value: source.PollingSchedule},
			{Key: "Config", Value: string(cfgJSON)},
			{Key: "Last Fetched", Value: source.LastFetchedAt},
			{Key: "Created", Value: source.CreatedAt},
		}
		if source.Prompt != "" {
			fields = append(fields, output.Field{Key: "Prompt", Value: source.Prompt})
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var sourcesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a source",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		sourceType, _ := cmd.Flags().GetString("type")
		configStr, _ := cmd.Flags().GetString("config")
		urlStr, _ := cmd.Flags().GetString("url")
		prompt, _ := cmd.Flags().GetString("prompt")

		if name == "" || sourceType == "" || (configStr == "" && urlStr == "") {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Source name").Value(&name),
					huh.NewInput().Title("Type (e.g. Sources::RssFeed or Sources::ZillowListing)").Value(&sourceType),
					huh.NewInput().Title("Config JSON").Description("Optional when using --url for Zillow sources").Value(&configStr),
					huh.NewInput().Title("Zillow URL").Description("Used to build {\"url\": \"...\"} when config JSON is omitted").Value(&urlStr),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		var cfg map[string]any
		if configStr != "" {
			if err := json.Unmarshal([]byte(configStr), &cfg); err != nil {
				return fmt.Errorf("invalid config JSON: %w", err)
			}
		}
		if urlStr != "" {
			if cfg == nil {
				cfg = map[string]any{}
			}
			cfg["url"] = urlStr
		}

		client := mustClient()
		svc := api.NewSourceService(client)
		source, err := svc.Create(cmdContext(), name, sourceType, "", prompt, true, cfg)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", source.ID)},
			{Key: "Name", Value: source.Name},
			{Key: "Type", Value: source.Type},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var sourcesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source ID: %s", args[0])
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
		if cmd.Flags().Changed("prompt") {
			v, _ := cmd.Flags().GetString("prompt")
			fields["prompt"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --name, --active, --prompt)")
		}

		client := mustClient()
		svc := api.NewSourceService(client)
		source, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", source.ID)},
			{Key: "Name", Value: source.Name},
			{Key: "Active", Value: fmt.Sprintf("%t", source.Active)},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var sourcesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source ID: %s", args[0])
		}

		var confirm bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Delete source %d?", id)).
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
		svc := api.NewSourceService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Source deleted.")
		return nil
	},
}

var sourcesFetchCmd = &cobra.Command{
	Use:   "fetch <id>",
	Short: "Trigger an async fetch for a source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewSourceService(client)
		_, msg, err := svc.Fetch(cmdContext(), id)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sourcesCmd)

	sourcesCmd.AddCommand(sourcesListCmd)
	sourcesListCmd.Flags().String("active", "", "filter by active status (true/false)")
	var page, perPage int
	addPaginationFlags(sourcesListCmd, &page, &perPage)

	sourcesCmd.AddCommand(sourcesShowCmd)

	sourcesCmd.AddCommand(sourcesCreateCmd)
	sourcesCreateCmd.Flags().String("name", "", "source name")
	sourcesCreateCmd.Flags().String("type", "", "source type (e.g. Sources::RssFeed or Sources::ZillowListing)")
	sourcesCreateCmd.Flags().String("config", "", "source config as JSON")
	sourcesCreateCmd.Flags().String("url", "", "convenience source URL; for Zillow this builds config {\"url\": \"...\"}")
	sourcesCreateCmd.Flags().String("prompt", "", "processing prompt")

	sourcesCmd.AddCommand(sourcesUpdateCmd)
	sourcesUpdateCmd.Flags().String("name", "", "new name")
	sourcesUpdateCmd.Flags().Bool("active", true, "active status")
	sourcesUpdateCmd.Flags().String("prompt", "", "new prompt")

	sourcesCmd.AddCommand(sourcesDeleteCmd)
	sourcesCmd.AddCommand(sourcesFetchCmd)
}
