package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Manage media sessions",
}

var mediaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List media sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewMediaSessionService(client)

		mediaType, _ := cmd.Flags().GetString("type")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.MediaListParams{
			Type: mediaType, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Title", "Created"}
		rows := make([][]string, len(resp.Items))
		for i, m := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", m.ID), m.Title, m.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var mediaShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show media session with turns",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		detail, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		m := detail.MediaSession
		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", m.ID)},
			{Key: "Title", Value: m.Title},
			{Key: "Created", Value: m.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(detail.MediaTurns) > 0 {
			fmt.Println()
			headers := []string{"ID", "Type", "Position", "Prompt", "Status", "Aspect Ratio"}
			rows := make([][]string, len(detail.MediaTurns))
			for i, t := range detail.MediaTurns {
				prompt := t.UserPrompt
				if len(prompt) > 50 {
					prompt = prompt[:50] + "..."
				}
				rows[i] = []string{
					fmt.Sprintf("%d", t.ID), t.Type,
					fmt.Sprintf("%d", t.Position), prompt,
					t.Status, t.AspectRatio,
				}
			}
			fmt.Print(formatter.FormatList(headers, rows, nil))
		}

		return nil
	},
}

var mediaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a media session",
	Long:  "Create a media session with an initial turn.\nExamples:\n  content_engine media create --type image --prompt \"A sunset\"\n  content_engine media create --type video --prompt \"Ocean waves\" --duration 8\n  content_engine media create --type audio --prompt \"Lo-fi beat\" --seconds 30",
	RunE: func(cmd *cobra.Command, args []string) error {
		mediaType, _ := cmd.Flags().GetString("type")
		prompt, _ := cmd.Flags().GetString("prompt")
		aspectRatio, _ := cmd.Flags().GetString("aspect-ratio")
		duration, _ := cmd.Flags().GetInt("duration")
		seconds, _ := cmd.Flags().GetInt("seconds")

		if mediaType == "" || prompt == "" {
			if !isInteractiveTerminal() {
				if mediaType == "" {
					return fmt.Errorf("--type is required")
				}
				return fmt.Errorf("--prompt is required")
			}
			var typeStr string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Media type").
						Options(
							huh.NewOption("Image", "image"),
							huh.NewOption("Video", "video"),
							huh.NewOption("Audio", "audio"),
						).
						Value(&typeStr),
					huh.NewText().Title("Prompt").Value(&prompt),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
			if mediaType == "" {
				mediaType = typeStr
			}
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		session, err := svc.Create(cmdContext(), mediaType, prompt, aspectRatio, duration, seconds)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", session.ID)},
			{Key: "Title", Value: session.Title},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Media generation started asynchronously.")
		return nil
	},
}

var mediaDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a media session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}

		confirm, err := confirmDestructiveAction(cmd, fmt.Sprintf("Delete media session %d?", id), "This will delete all turns and generated media.")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Media session deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mediaCmd)

	mediaCmd.AddCommand(mediaListCmd)
	mediaListCmd.Flags().String("type", "", "filter by type (image, video, audio)")
	var page, perPage int
	addPaginationFlags(mediaListCmd, &page, &perPage)

	mediaCmd.AddCommand(mediaShowCmd)

	mediaCmd.AddCommand(mediaCreateCmd)
	mediaCreateCmd.Flags().String("type", "", "media type: image, video, audio")
	mediaCreateCmd.Flags().String("prompt", "", "generation prompt")
	mediaCreateCmd.Flags().String("aspect-ratio", "", "aspect ratio (e.g. 16:9)")
	mediaCreateCmd.Flags().Int("duration", 0, "video duration in seconds")
	mediaCreateCmd.Flags().Int("seconds", 0, "audio capture seconds")

	mediaCmd.AddCommand(mediaDeleteCmd)
	addAutoConfirmFlags(mediaDeleteCmd)
}
