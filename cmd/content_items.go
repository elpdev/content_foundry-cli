package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var contentItemsCmd = &cobra.Command{
	Use:     "content-items",
	Aliases: []string{"ci"},
	Short:   "Manage content items",
}

var contentItemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List content items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewContentItemService(client)

		status, _ := cmd.Flags().GetString("status")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.ContentItemListParams{
			Status: status, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Title", "Status", "Source ID", "Fetched At"}
		rows := make([][]string, len(resp.Items))
		for i, item := range resp.Items {
			title := item.Title
			if len(title) > 60 {
				title = title[:60] + "..."
			}
			rows[i] = []string{
				fmt.Sprintf("%d", item.ID), title, item.Status,
				fmt.Sprintf("%d", item.SourceID), item.FetchedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var contentItemsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show content item details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid content item ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewContentItemService(client)
		item, drafts, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", item.ID)},
			{Key: "Title", Value: item.Title},
			{Key: "Status", Value: item.Status},
			{Key: "Source ID", Value: fmt.Sprintf("%d", item.SourceID)},
			{Key: "Source URL", Value: item.SourceURL},
			{Key: "Fetched At", Value: item.FetchedAt},
			{Key: "Created", Value: item.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(drafts) > 0 {
			fmt.Println()
			headers := []string{"ID", "Platform ID", "Status", "Version", "Scheduled For"}
			rows := make([][]string, len(drafts))
			for i, d := range drafts {
				rows[i] = []string{
					fmt.Sprintf("%d", d.ID),
					fmt.Sprintf("%d", d.PlatformID),
					d.Status,
					fmt.Sprintf("%d", d.Version),
					d.ScheduledFor,
				}
			}
			fmt.Print(formatter.FormatList(headers, rows, nil))
		}

		return nil
	},
}

var contentItemsProcessCmd = &cobra.Command{
	Use:   "process <id>",
	Short: "Process a content item (async)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid content item ID: %s", args[0])
		}

		guidance, _ := cmd.Flags().GetString("guidance")

		client := mustClient()
		svc := api.NewContentItemService(client)
		msg, err := svc.Process(cmdContext(), id, guidance)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

var contentItemsGenerateDraftsCmd = &cobra.Command{
	Use:   "generate-drafts <id>",
	Short: "Generate missing drafts for a content item (async)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid content item ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewContentItemService(client)
		msg, err := svc.GenerateDrafts(cmdContext(), id)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(contentItemsCmd)

	contentItemsCmd.AddCommand(contentItemsListCmd)
	contentItemsListCmd.Flags().String("status", "", "filter by status (pending, processed, ...)")
	var page, perPage int
	addPaginationFlags(contentItemsListCmd, &page, &perPage)

	contentItemsCmd.AddCommand(contentItemsShowCmd)

	contentItemsCmd.AddCommand(contentItemsProcessCmd)
	contentItemsProcessCmd.Flags().String("guidance", "", "processing guidance")

	contentItemsCmd.AddCommand(contentItemsGenerateDraftsCmd)
}
