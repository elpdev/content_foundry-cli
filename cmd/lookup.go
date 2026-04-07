package cmd

import (
	"fmt"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/spf13/cobra"
)

var lookupCmd = &cobra.Command{
	Use:   "lookup <url>",
	Short: "Look up a publication or source content item by URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		queryURL := args[0]
		client := mustClient()
		ctx := cmdContext()

		publicationsSvc := api.NewPublicationService(client)
		publications, err := publicationsSvc.List(ctx, api.PublicationListParams{
			URL: queryURL, Page: 1, PerPage: 1,
		})
		if err != nil {
			return err
		}

		if len(publications.Items) > 0 {
			detail, err := publicationsSvc.Get(ctx, publications.Items[0].ID)
			if err != nil {
				return err
			}
			printPublicationDetail(detail)
			return nil
		}

		contentItemsSvc := api.NewContentItemService(client)
		contentItems, err := contentItemsSvc.List(ctx, api.ContentItemListParams{
			Search: queryURL, Page: 1, PerPage: 20,
		})
		if err != nil {
			return err
		}

		if len(contentItems.Items) == 0 {
			fmt.Println("No matches found")
			return nil
		}

		headers := []string{"ID", "Title", "Status", "Source ID", "Source URL", "Fetched At"}
		rows := make([][]string, len(contentItems.Items))
		for i, item := range contentItems.Items {
			title := item.Title
			sourceURL := item.SourceURL
			if outFormat == "table" {
				title = truncate(title, 60)
				sourceURL = truncate(sourceURL, 80)
			}
			rows[i] = []string{
				fmt.Sprintf("%d", item.ID),
				title,
				item.Status,
				fmt.Sprintf("%d", item.SourceID),
				sourceURL,
				item.FetchedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &contentItems.Pagination))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
