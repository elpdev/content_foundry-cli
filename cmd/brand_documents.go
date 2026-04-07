package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Manage brand documents",
}

var docsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List brand documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewBrandDocumentService(client)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.BrandDocumentListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Title", "Source", "Status"}
		rows := make([][]string, len(resp.Items))
		for i, d := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", d.ID), d.Title, d.SourceType, d.Status,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var docsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show document details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid document ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewBrandDocumentService(client)
		doc, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", doc.ID)},
			{Key: "Title", Value: doc.Title},
			{Key: "Source Type", Value: doc.SourceType},
			{Key: "Source URL", Value: doc.SourceURL},
			{Key: "Status", Value: doc.Status},
			{Key: "Created", Value: doc.CreatedAt},
		}
		if doc.ErrorMessage != "" {
			fields = append(fields, output.Field{Key: "Error", Value: doc.ErrorMessage})
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var docsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a document from URL or file",
	RunE: func(cmd *cobra.Command, args []string) error {
		docURL, _ := cmd.Flags().GetString("url")

		if docURL == "" {
			return fmt.Errorf("--url is required")
		}

		client := mustClient()
		svc := api.NewBrandDocumentService(client)
		doc, msg, err := svc.CreateFromURL(cmdContext(), docURL)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", doc.ID)},
			{Key: "Title", Value: doc.Title},
			{Key: "Status", Value: doc.Status},
		}
		fmt.Print(formatter.FormatItem(fields))
		if msg != "" {
			fmt.Println(msg)
		}
		return nil
	},
}

var docsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid document ID: %s", args[0])
		}

		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Delete document %d?", id)).
			Description("This cannot be undone.").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewBrandDocumentService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Document deleted.")
		return nil
	},
}

var docsIndexContentCmd = &cobra.Command{
	Use:   "index-content",
	Short: "Trigger content indexing for all documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewBrandDocumentService(client)
		msg, err := svc.IndexContent(cmdContext())
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.AddCommand(docsListCmd)
	var page, perPage int
	addPaginationFlags(docsListCmd, &page, &perPage)

	docsCmd.AddCommand(docsShowCmd)

	docsCmd.AddCommand(docsCreateCmd)
	docsCreateCmd.Flags().String("url", "", "URL to ingest")

	docsCmd.AddCommand(docsDeleteCmd)
	docsCmd.AddCommand(docsIndexContentCmd)
}
