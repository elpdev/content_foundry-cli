package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/models"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var publicationsCmd = &cobra.Command{
	Use:   "publications",
	Short: "Manage publications",
}

var publicationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewPublicationService(client)

		urlFilter, _ := cmd.Flags().GetString("url")
		status, _ := cmd.Flags().GetString("status")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.PublicationListParams{
			URL: urlFilter, Status: status, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Draft", "Status", "URL", "Content Item", "Platform", "Published At"}
		rows := make([][]string, len(resp.Items))
		for i, pub := range resp.Items {
			urlValue := pub.URL
			title := pub.ContentItemTitle
			if outFormat == "table" {
				urlValue = truncate(urlValue, 80)
				title = truncate(title, 60)
			}
			rows[i] = []string{
				fmt.Sprintf("%d", pub.ID),
				fmt.Sprintf("%d", pub.DraftID),
				pub.Status,
				urlValue,
				title,
				fmt.Sprintf("%d", pub.PlatformID),
				pub.PublishedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var publicationsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show publication details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid publication ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewPublicationService(client)
		detail, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		printPublicationDetail(detail)
		return nil
	},
}

func printPublicationDetail(detail *models.PublicationDetail) {
	pub := detail.Publication
	draft := detail.Draft
	item := detail.ContentItem
	source := detail.Source

	content := draft.Content
	if outFormat == "table" {
		content = truncate(content, 200)
	}

	fields := []output.Field{
		{Key: "Publication ID", Value: fmt.Sprintf("%d", pub.ID)},
		{Key: "Publication URL", Value: pub.URL},
		{Key: "Publication Status", Value: pub.Status},
		{Key: "Published At", Value: pub.PublishedAt},
		{Key: "Draft ID", Value: fmt.Sprintf("%d", draft.ID)},
		{Key: "Draft Status", Value: draft.Status},
		{Key: "Draft Content", Value: content},
		{Key: "Content Item ID", Value: fmt.Sprintf("%d", item.ID)},
		{Key: "Content Item Title", Value: item.Title},
		{Key: "Source URL", Value: item.SourceURL},
		{Key: "Source ID", Value: fmt.Sprintf("%d", source.ID)},
		{Key: "Source Name", Value: source.Name},
		{Key: "Source Type", Value: source.Type},
	}
	if pub.ErrorMessage != "" {
		fields = append(fields, output.Field{Key: "Error", Value: pub.ErrorMessage})
	}

	fmt.Print(formatter.FormatItem(fields))
}

func init() {
	rootCmd.AddCommand(publicationsCmd)

	publicationsCmd.AddCommand(publicationsListCmd)
	publicationsListCmd.Flags().String("url", "", "filter by publication URL")
	publicationsListCmd.Flags().String("status", "", "filter by status (pending, publishing, published, failed)")
	var page, perPage int
	addPaginationFlags(publicationsListCmd, &page, &perPage)

	publicationsCmd.AddCommand(publicationsShowCmd)
}
