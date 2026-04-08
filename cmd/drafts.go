package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/models"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var draftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "Manage drafts and workflow",
}

var draftsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewDraftService(client)

		status, _ := cmd.Flags().GetString("status")
		platformID, _ := cmd.Flags().GetString("platform-id")
		assignedTo, _ := cmd.Flags().GetString("assigned-to")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.DraftListParams{
			Status: status, PlatformID: platformID, AssignedToID: assignedTo,
			Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Content Item", "Platform", "Status", "Version", "Scheduled For"}
		rows := make([][]string, len(resp.Items))
		for i, d := range resp.Items {
			rows[i] = []string{
				fmt.Sprintf("%d", d.ID),
				fmt.Sprintf("%d", d.ContentItemID),
				fmt.Sprintf("%d", d.PlatformID),
				d.Status,
				fmt.Sprintf("%d", d.Version),
				d.ScheduledFor,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var draftsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show draft details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		detail, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		d := detail.Draft
		content := d.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", d.ID)},
			{Key: "Content Item", Value: fmt.Sprintf("%d", d.ContentItemID)},
			{Key: "Platform", Value: fmt.Sprintf("%d", d.PlatformID)},
			{Key: "Status", Value: d.Status},
			{Key: "Version", Value: fmt.Sprintf("%d", d.Version)},
			{Key: "Content", Value: content},
			{Key: "Scheduled For", Value: d.ScheduledFor},
			{Key: "Created", Value: d.CreatedAt},
		}
		if d.AssignedToID != nil {
			fields = append(fields, output.Field{Key: "Assigned To", Value: fmt.Sprintf("%d", *d.AssignedToID)})
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(detail.Comments) > 0 {
			fmt.Println()
			fmt.Println("Comments:")
			headers := []string{"ID", "User", "Body", "Created"}
			rows := make([][]string, len(detail.Comments))
			for i, c := range detail.Comments {
				body := c.Body
				if len(body) > 80 {
					body = body[:80] + "..."
				}
				rows[i] = []string{
					fmt.Sprintf("%d", c.ID),
					fmt.Sprintf("%d", c.UserID),
					body,
					c.CreatedAt,
				}
			}
			fmt.Print(formatter.FormatList(headers, rows, nil))
		}

		if detail.Publication != nil {
			fmt.Println()
			fmt.Println("Publication:")
			pub := detail.Publication
			pubFields := []output.Field{
				{Key: "ID", Value: fmt.Sprintf("%d", pub.ID)},
				{Key: "Status", Value: pub.Status},
				{Key: "URL", Value: pub.URL},
				{Key: "Published At", Value: pub.PublishedAt},
			}
			if pub.ErrorMessage != "" {
				pubFields = append(pubFields, output.Field{Key: "Error", Value: pub.ErrorMessage})
			}
			fmt.Print(formatter.FormatItem(pubFields))
		}

		return nil
	},
}

var draftsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a draft",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		platformID, _ := cmd.Flags().GetInt64("platform-id")

		if title == "" || content == "" || platformID == 0 {
			if !isInteractiveTerminal() {
				switch {
				case title == "":
					return fmt.Errorf("--title is required")
				case content == "":
					return fmt.Errorf("--content is required")
				default:
					return fmt.Errorf("--platform-id is required")
				}
			}
			var pidStr string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Title").Value(&title),
					huh.NewText().Title("Content").Value(&content),
					huh.NewInput().Title("Platform ID").Value(&pidStr),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
			if pid, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
				platformID = pid
			}
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Create(cmdContext(), title, content, platformID)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", draft.ID)},
			{Key: "Status", Value: draft.Status},
			{Key: "Version", Value: fmt.Sprintf("%d", draft.Version)},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var draftsEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a draft's title or content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("title") {
			value, _ := cmd.Flags().GetString("title")
			fields["title"] = value
		}
		if cmd.Flags().Changed("content") {
			value, _ := cmd.Flags().GetString("content")
			fields["content"] = value
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --title and/or --content)")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		printDraftStatus(draft)
		return nil
	},
}

var draftsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		confirm, err := confirmDestructiveAction(cmd, fmt.Sprintf("Delete draft %d?", id), "This cannot be undone.")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Draft deleted.")
		return nil
	},
}

// Workflow commands

var draftsApproveCmd = &cobra.Command{
	Use:   "approve <id>",
	Short: "Approve a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Approve(cmdContext(), id)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsRejectCmd = &cobra.Command{
	Use:   "reject <id>",
	Short: "Reject a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Reject(cmdContext(), id)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsReviseCmd = &cobra.Command{
	Use:   "revise <id>",
	Short: "Request revision on a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		notes, _ := cmd.Flags().GetString("notes")
		if notes == "" {
			return fmt.Errorf("--notes is required")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.RequestRevision(cmdContext(), id, notes)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsScheduleCmd = &cobra.Command{
	Use:   "schedule <id>",
	Short: "Schedule a draft for publication",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		at, _ := cmd.Flags().GetString("at")
		if at == "" {
			return fmt.Errorf("--at is required (e.g. 2026-04-10T14:00:00Z)")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Schedule(cmdContext(), id, at)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsRescheduleCmd = &cobra.Command{
	Use:   "reschedule <id>",
	Short: "Reschedule a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		at, _ := cmd.Flags().GetString("at")
		if at == "" {
			return fmt.Errorf("--at is required")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Reschedule(cmdContext(), id, at)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsUnscheduleCmd = &cobra.Command{
	Use:   "unschedule <id>",
	Short: "Unschedule a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Unschedule(cmdContext(), id)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsAssignCmd = &cobra.Command{
	Use:   "assign <id>",
	Short: "Assign a draft to a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		to, _ := cmd.Flags().GetInt64("to")
		if to == 0 {
			return fmt.Errorf("--to <user_id> is required")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Assign(cmdContext(), id, to)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsUnassignCmd = &cobra.Command{
	Use:   "unassign <id>",
	Short: "Unassign a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.Unassign(cmdContext(), id)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsMediaCmd = &cobra.Command{
	Use:   "media <id>",
	Short: "Save media turns to a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		turnIDsStr, _ := cmd.Flags().GetString("turn-ids")
		if turnIDsStr == "" {
			return fmt.Errorf("--turn-ids is required (comma-separated)")
		}

		turnIDs, err := parseIntList(turnIDsStr)
		if err != nil {
			return fmt.Errorf("invalid turn IDs: %w", err)
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		draft, err := svc.SaveMedia(cmdContext(), id, turnIDs)
		if err != nil {
			return err
		}
		printDraftStatus(draft)
		return nil
	},
}

var draftsPublishCmd = &cobra.Command{
	Use:   "publish <id>",
	Short: "Publish a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}

		var turnIDs []int64
		turnIDsStr, _ := cmd.Flags().GetString("turn-ids")
		if turnIDsStr != "" {
			turnIDs, err = parseIntList(turnIDsStr)
			if err != nil {
				return fmt.Errorf("invalid turn IDs: %w", err)
			}
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		pub, msg, err := svc.Publish(cmdContext(), id, turnIDs)
		if err != nil {
			return err
		}

		if msg != "" {
			fmt.Println(msg)
		}
		fields := []output.Field{
			{Key: "Publication ID", Value: fmt.Sprintf("%d", pub.ID)},
			{Key: "Status", Value: pub.Status},
		}
		if pub.URL != "" {
			fields = append(fields, output.Field{Key: "URL", Value: pub.URL})
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

// Comments subcommands

var draftsCommentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Manage draft comments",
}

var draftsCommentsAddCmd = &cobra.Command{
	Use:   "add <draft_id>",
	Short: "Add a comment to a draft",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		draftID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			return fmt.Errorf("--body is required")
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		comment, err := svc.AddComment(cmdContext(), draftID, body)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", comment.ID)},
			{Key: "Body", Value: comment.Body},
			{Key: "Created", Value: comment.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var draftsCommentsDeleteCmd = &cobra.Command{
	Use:   "delete <draft_id> <comment_id>",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		draftID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid draft ID: %s", args[0])
		}
		commentID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid comment ID: %s", args[1])
		}

		confirm, err := confirmDestructiveAction(cmd, fmt.Sprintf("Delete comment %d?", commentID), "This cannot be undone.")
		if err != nil {
			return err
		}
		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewDraftService(client)
		if err := svc.DeleteComment(cmdContext(), draftID, commentID); err != nil {
			return err
		}
		fmt.Println("Comment deleted.")
		return nil
	},
}

// Helpers

func printDraftStatus(d *models.Draft) {
	fields := []output.Field{
		{Key: "ID", Value: fmt.Sprintf("%d", d.ID)},
		{Key: "Status", Value: d.Status},
		{Key: "Version", Value: fmt.Sprintf("%d", d.Version)},
	}
	if d.ScheduledFor != "" {
		fields = append(fields, output.Field{Key: "Scheduled For", Value: d.ScheduledFor})
	}
	if d.AssignedToID != nil {
		fields = append(fields, output.Field{Key: "Assigned To", Value: fmt.Sprintf("%d", *d.AssignedToID)})
	}
	fmt.Print(formatter.FormatItem(fields))
}

func parseIntList(s string) ([]int64, error) {
	parts := strings.Split(s, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ID %q", p)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func init() {
	rootCmd.AddCommand(draftsCmd)

	draftsCmd.AddCommand(draftsListCmd)
	draftsListCmd.Flags().String("status", "", "filter by status (pending_review, approved, rejected, scheduled, published)")
	draftsListCmd.Flags().String("platform-id", "", "filter by platform ID")
	draftsListCmd.Flags().String("assigned-to", "", "filter by assigned user ID")
	var page, perPage int
	addPaginationFlags(draftsListCmd, &page, &perPage)

	draftsCmd.AddCommand(draftsShowCmd)

	draftsCmd.AddCommand(draftsCreateCmd)
	draftsCreateCmd.Flags().String("title", "", "draft title")
	draftsCreateCmd.Flags().String("content", "", "draft content")
	draftsCreateCmd.Flags().Int64("platform-id", 0, "platform ID")

	draftsCmd.AddCommand(draftsEditCmd)
	draftsEditCmd.Flags().String("title", "", "new draft title")
	draftsEditCmd.Flags().String("content", "", "new draft content")

	draftsCmd.AddCommand(draftsDeleteCmd)
	addAutoConfirmFlags(draftsDeleteCmd)
	draftsCmd.AddCommand(draftsApproveCmd)
	draftsCmd.AddCommand(draftsRejectCmd)

	draftsCmd.AddCommand(draftsReviseCmd)
	draftsReviseCmd.Flags().String("notes", "", "revision notes")

	draftsCmd.AddCommand(draftsScheduleCmd)
	draftsScheduleCmd.Flags().String("at", "", "scheduled time (ISO 8601)")

	draftsCmd.AddCommand(draftsRescheduleCmd)
	draftsRescheduleCmd.Flags().String("at", "", "new scheduled time (ISO 8601)")

	draftsCmd.AddCommand(draftsUnscheduleCmd)

	draftsCmd.AddCommand(draftsAssignCmd)
	draftsAssignCmd.Flags().Int64("to", 0, "user ID to assign to")

	draftsCmd.AddCommand(draftsUnassignCmd)

	draftsCmd.AddCommand(draftsMediaCmd)
	draftsMediaCmd.Flags().String("turn-ids", "", "comma-separated media turn IDs")

	draftsCmd.AddCommand(draftsPublishCmd)
	draftsPublishCmd.Flags().String("turn-ids", "", "comma-separated media turn IDs (optional)")

	draftsCmd.AddCommand(draftsCommentsCmd)
	draftsCommentsCmd.AddCommand(draftsCommentsAddCmd)
	draftsCommentsAddCmd.Flags().String("body", "", "comment body")
	draftsCommentsCmd.AddCommand(draftsCommentsDeleteCmd)
	addAutoConfirmFlags(draftsCommentsDeleteCmd)
}
