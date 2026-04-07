package cmd

import (
	"fmt"

	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/spf13/cobra"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "View activity events",
}

var activityListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent activity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewActivityService(client)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.ActivityListParams{
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Event", "Trackable", "User", "Created"}
		rows := make([][]string, len(resp.Items))
		for i, e := range resp.Items {
			trackable := e.TrackableType
			if e.TrackableID > 0 {
				trackable += fmt.Sprintf(" #%d", e.TrackableID)
			}
			rows[i] = []string{
				fmt.Sprintf("%d", e.ID), e.EventType, trackable,
				fmt.Sprintf("%d", e.UserID), e.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(activityCmd)

	activityCmd.AddCommand(activityListCmd)
	var page, perPage int
	addPaginationFlags(activityListCmd, &page, &perPage)
}
