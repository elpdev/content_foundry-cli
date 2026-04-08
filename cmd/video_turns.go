package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var videosCmd = &cobra.Command{
	Use:   "videos",
	Short: "Manage video turns",
}

var videosShowCmd = &cobra.Command{
	Use:   "show <media_session_id> <turn_id>",
	Short: "Show a video turn",
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
		detail, err := svc.Get(cmdContext(), sessionID)
		if err != nil {
			return err
		}

		for _, t := range detail.MediaTurns {
			if t.ID == turnID {
				fields := []output.Field{
					{Key: "ID", Value: fmt.Sprintf("%d", t.ID)},
					{Key: "Session", Value: fmt.Sprintf("%d", t.MediaSessionID)},
					{Key: "Position", Value: fmt.Sprintf("%d", t.Position)},
					{Key: "Prompt", Value: t.UserPrompt},
					{Key: "Status", Value: t.Status},
					{Key: "Aspect Ratio", Value: t.AspectRatio},
					{Key: "Duration", Value: fmt.Sprintf("%ds", t.DurationSeconds)},
					{Key: "Created", Value: t.CreatedAt},
				}
				if t.VideoURL != "" {
					fields = append(fields, output.Field{Key: "URL", Value: t.VideoURL})
				}
				fmt.Print(formatter.FormatItem(fields))
				return nil
			}
		}

		return fmt.Errorf("turn %d not found in media session %d", turnID, sessionID)
	},
}

var videosCreateCmd = &cobra.Command{
	Use:   "create <media_session_id>",
	Short: "Create a new video turn",
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
		duration, _ := cmd.Flags().GetInt("duration")

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		turn, err := svc.CreateVideoTurn(cmdContext(), sessionID, prompt, aspectRatio, duration)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", turn.ID)},
			{Key: "Status", Value: turn.Status},
			{Key: "Prompt", Value: turn.UserPrompt},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Video generation started asynchronously.")
		return nil
	},
}

var videosExtendCmd = &cobra.Command{
	Use:   "extend <media_session_id> <turn_id>",
	Short: "Extend a video turn",
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

		targetDuration, _ := cmd.Flags().GetInt("target-duration")
		if targetDuration == 0 {
			return fmt.Errorf("--target-duration is required")
		}
		extensionPrompt, _ := cmd.Flags().GetString("prompt")

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		turn, err := svc.ExtendVideoTurn(cmdContext(), sessionID, turnID, targetDuration, extensionPrompt)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", turn.ID)},
			{Key: "Status", Value: turn.Status},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Video extension started asynchronously.")
		return nil
	},
}

var videosDeleteCmd = &cobra.Command{
	Use:   "delete <media_session_id> <turn_id>",
	Short: "Delete a video turn",
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

		confirm, err := confirmDestructiveAction(cmd, fmt.Sprintf("Delete video turn %d?", turnID), "This cannot be undone.")
		if err != nil {
			return err
		}
		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		if err := svc.DeleteVideoTurn(cmdContext(), sessionID, turnID); err != nil {
			return err
		}
		fmt.Println("Video turn deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(videosCmd)

	videosCmd.AddCommand(videosShowCmd)

	videosCmd.AddCommand(videosCreateCmd)
	videosCreateCmd.Flags().String("prompt", "", "video prompt")
	videosCreateCmd.Flags().String("aspect-ratio", "", "aspect ratio (16:9 or 9:16)")
	videosCreateCmd.Flags().Int("duration", 0, "duration in seconds (4, 6, or 8)")

	videosCmd.AddCommand(videosExtendCmd)
	videosExtendCmd.Flags().Int("target-duration", 0, "target total duration in seconds")
	videosExtendCmd.Flags().String("prompt", "", "extension prompt (optional)")

	videosCmd.AddCommand(videosDeleteCmd)
	addAutoConfirmFlags(videosDeleteCmd)
}
