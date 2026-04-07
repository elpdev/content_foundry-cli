package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Manage audio turns",
}

var audioCreateCmd = &cobra.Command{
	Use:   "create <media_session_id>",
	Short: "Create a new audio turn",
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
		seconds, _ := cmd.Flags().GetInt("seconds")

		client := mustClient()
		svc := api.NewMediaSessionService(client)
		turn, err := svc.CreateAudioTurn(cmdContext(), sessionID, prompt, seconds)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", turn.ID)},
			{Key: "Status", Value: turn.Status},
			{Key: "Prompt", Value: turn.UserPrompt},
			{Key: "Duration", Value: fmt.Sprintf("%ds", turn.CaptureSeconds)},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Audio generation started asynchronously.")
		return nil
	},
}

var audioDeleteCmd = &cobra.Command{
	Use:   "delete <media_session_id> <turn_id>",
	Short: "Delete an audio turn",
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
		if err := svc.DeleteAudioTurn(cmdContext(), sessionID, turnID); err != nil {
			return err
		}
		fmt.Println("Audio turn deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(audioCmd)

	audioCmd.AddCommand(audioCreateCmd)
	audioCreateCmd.Flags().String("prompt", "", "audio prompt")
	audioCreateCmd.Flags().Int("seconds", 30, "capture duration in seconds")

	audioCmd.AddCommand(audioDeleteCmd)
}
