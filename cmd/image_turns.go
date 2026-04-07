package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/media"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Manage image turns",
}

var imagesShowCmd = &cobra.Command{
	Use:   "show <media_session_id> <turn_id>",
	Short: "Show an image turn (renders inline if supported)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}
		turnID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid turn ID: %s", args[1])
		}

		// Fetch the media session to find this turn
		client := mustClient()
		svc := api.NewMediaSessionService(client)
		detail, err := svc.Get(cmdContext(), sessionID)
		if err != nil {
			return err
		}

		// Find the image turn in the session's turns
		var found bool
		for _, t := range detail.MediaTurns {
			if t.ID == turnID {
				found = true
				fields := []output.Field{
					{Key: "ID", Value: fmt.Sprintf("%d", t.ID)},
					{Key: "Session", Value: fmt.Sprintf("%d", t.MediaSessionID)},
					{Key: "Position", Value: fmt.Sprintf("%d", t.Position)},
					{Key: "Prompt", Value: t.UserPrompt},
					{Key: "Status", Value: t.Status},
					{Key: "Aspect Ratio", Value: t.AspectRatio},
					{Key: "Created", Value: t.CreatedAt},
				}
				if t.ImageURL != "" {
					fields = append(fields, output.Field{Key: "URL", Value: t.ImageURL})
				}
				fmt.Print(formatter.FormatItem(fields))

				// Render image inline if terminal supports it
				if t.ImageURL != "" {
					rendered := media.RenderImageFromURL(t.ImageURL, 80)
					if rendered != "" {
						fmt.Println()
						fmt.Print(rendered)
						fmt.Println()
					}
				}
				break
			}
		}

		if !found {
			return fmt.Errorf("turn %d not found in media session %d", turnID, sessionID)
		}

		return nil
	},
}

var imagesCreateCmd = &cobra.Command{
	Use:   "create <media_session_id>",
	Short: "Create a new image turn",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}

		prompt, _ := cmd.Flags().GetString("prompt")
		if prompt == "" {
			return fmt.Errorf("--prompt is required")
		}
		aspectRatio, _ := cmd.Flags().GetString("aspect-ratio")

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		turn, err := svc.CreateImageTurn(cmdContext(), sessionID, prompt, aspectRatio)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", turn.ID)},
			{Key: "Status", Value: turn.Status},
			{Key: "Prompt", Value: turn.UserPrompt},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Image generation started asynchronously.")
		return nil
	},
}

var imagesConvertCmd = &cobra.Command{
	Use:   "convert <media_session_id> <turn_id>",
	Short: "Convert an image to multiple formats",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}
		turnID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid turn ID: %s", args[1])
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		turn, err := svc.ConvertImageTurn(cmdContext(), sessionID, turnID)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", turn.ID)},
			{Key: "Status", Value: turn.Status},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Image conversion started.")
		return nil
	},
}

var imagesDeleteCmd = &cobra.Command{
	Use:   "delete <media_session_id> <turn_id>",
	Short: "Delete an image turn",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid media session ID: %s", args[0])
		}
		turnID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid turn ID: %s", args[1])
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		if err := svc.DeleteImageTurn(cmdContext(), sessionID, turnID); err != nil {
			return err
		}
		fmt.Println("Image turn deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)

	imagesCmd.AddCommand(imagesShowCmd)

	imagesCmd.AddCommand(imagesCreateCmd)
	imagesCreateCmd.Flags().String("prompt", "", "image prompt")
	imagesCreateCmd.Flags().String("aspect-ratio", "", "aspect ratio (e.g. 16:9, 1:1)")

	imagesCmd.AddCommand(imagesConvertCmd)
	imagesCmd.AddCommand(imagesDeleteCmd)
}
