package cmd

import (
	"fmt"
	"strconv"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/spf13/cobra"
)

var notificationsCmd = &cobra.Command{
	Use:     "notifications",
	Aliases: []string{"notif"},
	Short:   "Manage notifications",
}

var notificationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewNotificationService(client)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.NotificationListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Type", "Read", "Created"}
		rows := make([][]string, len(resp.Items))
		for i, n := range resp.Items {
			read := "-"
			if n.ReadAt != "" {
				read = n.ReadAt
			}
			rows[i] = []string{
				fmt.Sprintf("%d", n.ID), n.Type, read, n.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var notificationsReadCmd = &cobra.Command{
	Use:   "read <id>",
	Short: "Mark a notification as read",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid notification ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewNotificationService(client)
		if _, err := svc.MarkRead(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Notification marked as read.")
		return nil
	},
}

var notificationsReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all notifications as read",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewNotificationService(client)
		msg, err := svc.MarkAllRead(cmdContext())
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)

	notificationsCmd.AddCommand(notificationsListCmd)
	var page, perPage int
	addPaginationFlags(notificationsListCmd, &page, &perPage)

	notificationsCmd.AddCommand(notificationsReadCmd)
	notificationsCmd.AddCommand(notificationsReadAllCmd)
}
