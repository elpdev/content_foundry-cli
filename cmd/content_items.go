package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/models"
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
		search, _ := cmd.Flags().GetString("search")
		sourceType, _ := cmd.Flags().GetString("source-type")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.ContentItemListParams{
			Status: status, Search: search, SourceType: sourceType, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Title", "Status", "Source ID", "Images", "Fetched At"}
		if sourceType == "inbound_email" {
			headers = []string{"ID", "Title", "Status", "Source ID", "From", "Received At", "Attachments", "Images", "Fetched At"}
		}
		rows := make([][]string, len(resp.Items))
		for i, item := range resp.Items {
			title := item.Title
			if outFormat == "table" {
				title = truncate(title, 60)
			}
			if sourceType == "inbound_email" {
				rows[i] = []string{
					fmt.Sprintf("%d", item.ID), title, item.Status,
					fmt.Sprintf("%d", item.SourceID), contentItemEmailFrom(&item), contentItemEmailReceivedAt(&item), contentItemEmailAttachmentCount(&item), contentItemImageSummary(&item), item.FetchedAt,
				}
				continue
			}
			rows[i] = []string{
				fmt.Sprintf("%d", item.ID), title, item.Status,
				fmt.Sprintf("%d", item.SourceID), contentItemImageSummary(&item), item.FetchedAt,
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
		if outFormat == "json" {
			raw, err := svc.GetRaw(cmdContext(), id)
			if err != nil {
				return err
			}
			fmt.Print(formatter.FormatRaw(raw))
			return nil
		}

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
		if item.HeroImageURL != "" {
			fields = append(fields, output.Field{Key: "Hero Image", Value: item.HeroImageURL})
		}
		if from := contentItemEmailFrom(item); from != "" {
			fields = append(fields, output.Field{Key: "From", Value: from})
		}
		if receivedAt := contentItemEmailReceivedAt(item); receivedAt != "" {
			fields = append(fields, output.Field{Key: "Received At", Value: receivedAt})
		}
		if attachmentCount, ok := contentItemEmailAttachmentCountIfPresent(item); ok {
			fields = append(fields, output.Field{Key: "Attachment Count", Value: attachmentCount})
		}
		if len(item.Assets) > 0 {
			fields = append(fields, output.Field{Key: "Linked Assets", Value: fmt.Sprintf("%d", len(item.Assets))})
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(item.Assets) > 0 {
			fmt.Println()
			headers := []string{"ID", "Filename", "Content Type", "URL"}
			rows := make([][]string, len(item.Assets))
			for i, asset := range item.Assets {
				url := asset.FileURL
				if outFormat == "table" {
					url = truncate(url, 80)
				}
				rows[i] = []string{
					fmt.Sprintf("%d", asset.ID),
					asset.Filename,
					asset.ContentType,
					url,
				}
			}
			fmt.Print(formatter.FormatList(headers, rows, nil))
		}

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

func contentItemImageSummary(item *models.ContentItem) string {
	if len(item.Assets) > 0 {
		return fmt.Sprintf("%d", len(item.Assets))
	}
	if item.HeroImageURL != "" {
		return "1"
	}
	return "0"
}

func contentItemEmailFrom(item *models.ContentItem) string {
	return contentItemMetadataString(item, "from")
}

func contentItemEmailReceivedAt(item *models.ContentItem) string {
	return contentItemMetadataString(item, "received_at")
}

func contentItemEmailAttachmentCount(item *models.ContentItem) string {
	if count, ok := contentItemEmailAttachmentCountIfPresent(item); ok {
		return count
	}
	return "0"
}

func contentItemEmailAttachmentCountIfPresent(item *models.ContentItem) (string, bool) {
	if item == nil || item.Metadata == nil {
		return "", false
	}
	v, ok := item.Metadata["attachment_count"]
	if !ok || v == nil {
		return "", false
	}
	switch n := v.(type) {
	case string:
		return n, true
	case float64:
		return strconv.FormatInt(int64(n), 10), true
	case float32:
		return strconv.FormatInt(int64(n), 10), true
	case int:
		return strconv.Itoa(n), true
	case int64:
		return strconv.FormatInt(n, 10), true
	case int32:
		return strconv.FormatInt(int64(n), 10), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

func contentItemMetadataString(item *models.ContentItem, key string) string {
	if item == nil || item.Metadata == nil {
		return ""
	}
	v, ok := item.Metadata[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
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
	contentItemsListCmd.Flags().String("search", "", "search title or source URL")
	contentItemsListCmd.Flags().String("source-type", "", "filter by source type (e.g. inbound_email)")
	var page, perPage int
	addPaginationFlags(contentItemsListCmd, &page, &perPage)

	contentItemsCmd.AddCommand(contentItemsShowCmd)

	contentItemsCmd.AddCommand(contentItemsProcessCmd)
	contentItemsProcessCmd.Flags().String("guidance", "", "processing guidance")

	contentItemsCmd.AddCommand(contentItemsGenerateDraftsCmd)
}
