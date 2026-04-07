package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/media"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Manage uploaded assets",
}

var assetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAssetService(client)

		assetType, _ := cmd.Flags().GetString("type")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.AssetListParams{
			Type: assetType, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Filename", "Content Type", "Size", "Created"}
		rows := make([][]string, len(resp.Items))
		for i, a := range resp.Items {
			size := fmt.Sprintf("%d KB", a.ByteSize/1024)
			rows[i] = []string{
				fmt.Sprintf("%d", a.ID), a.Filename, a.ContentType, size, a.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var assetsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show asset details (renders inline for images)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid asset ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewAssetService(client)
		a, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", a.ID)},
			{Key: "Filename", Value: a.Filename},
			{Key: "Content Type", Value: a.ContentType},
			{Key: "Size", Value: fmt.Sprintf("%d bytes", a.ByteSize)},
			{Key: "Created", Value: a.CreatedAt},
		}
		if a.FileURL != "" {
			fields = append(fields, output.Field{Key: "URL", Value: a.FileURL})
		}
		fmt.Print(formatter.FormatItem(fields))

		// Render inline for image assets
		if a.FileURL != "" && strings.HasPrefix(a.ContentType, "image/") {
			rendered := media.RenderImageFromURL(a.FileURL, 80)
			if rendered != "" {
				fmt.Println()
				fmt.Print(rendered)
				fmt.Println()
			}
		}

		return nil
	},
}

var assetsUploadCmd = &cobra.Command{
	Use:   "upload <file_path>",
	Short: "Upload a file as an asset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewAssetService(client)

		asset, err := svc.Upload(cmdContext(), args[0])
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", asset.ID)},
			{Key: "Filename", Value: asset.Filename},
			{Key: "Content Type", Value: asset.ContentType},
			{Key: "Size", Value: fmt.Sprintf("%d bytes", asset.ByteSize)},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var assetsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an asset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid asset ID: %s", args[0])
		}

		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Delete asset %d?", id)).
			Description("This cannot be undone.").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewAssetService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Asset deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(assetsCmd)

	assetsCmd.AddCommand(assetsListCmd)
	assetsListCmd.Flags().String("type", "", "filter by type (image, video)")
	var page, perPage int
	addPaginationFlags(assetsListCmd, &page, &perPage)

	assetsCmd.AddCommand(assetsShowCmd)
	assetsCmd.AddCommand(assetsUploadCmd)
	assetsCmd.AddCommand(assetsDeleteCmd)
}
