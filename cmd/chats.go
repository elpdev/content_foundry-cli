package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Manage chat sessions",
}

var chatsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List chats",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewChatService(client)

		personaID, _ := cmd.Flags().GetString("persona-id")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		resp, err := svc.List(cmdContext(), api.ChatListParams{
			PersonaID: personaID, Page: page, PerPage: perPage,
		})
		if err != nil {
			return err
		}

		headers := []string{"ID", "Title", "Persona", "Model", "Created"}
		rows := make([][]string, len(resp.Items))
		for i, c := range resp.Items {
			title := c.Title
			if title == "" {
				title = "Untitled"
			}
			personaStr := "-"
			if c.PersonaID != nil {
				personaStr = fmt.Sprintf("%d", *c.PersonaID)
			}
			rows[i] = []string{
				fmt.Sprintf("%d", c.ID), title, personaStr,
				fmt.Sprintf("%v", c.ModelID), c.CreatedAt,
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, &resp.Pagination))
		return nil
	},
}

var chatsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show chat with messages",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewChatService(client)
		detail, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		c := detail.Chat
		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", c.ID)},
			{Key: "Title", Value: c.Title},
			{Key: "Model", Value: fmt.Sprintf("%v", c.ModelID)},
			{Key: "Created", Value: c.CreatedAt},
		}
		if c.PersonaID != nil {
			fields = append(fields, output.Field{Key: "Persona", Value: fmt.Sprintf("%d", *c.PersonaID)})
		}
		fmt.Print(formatter.FormatItem(fields))

		if len(detail.Messages) > 0 {
			fmt.Println()
			for _, m := range detail.Messages {
				content := m.Content
				if len(content) > 300 {
					content = content[:300] + "..."
				}
				fmt.Printf("[%s] %s\n\n", m.Role, content)
			}
		}

		return nil
	},
}

var chatsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt, _ := cmd.Flags().GetString("prompt")
		personaIDVal, _ := cmd.Flags().GetInt64("persona-id")
		model, _ := cmd.Flags().GetString("model")

		if prompt == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewText().Title("Initial prompt").Value(&prompt),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		var personaID *int64
		if personaIDVal > 0 {
			personaID = &personaIDVal
		}

		client := mustClient()
		svc := api.NewChatService(client)
		chat, err := svc.Create(cmdContext(), prompt, personaID, model)
		if err != nil {
			return err
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", chat.ID)},
			{Key: "Title", Value: chat.Title},
			{Key: "Created", Value: chat.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		fmt.Println("Chat created. Response is being generated asynchronously.")
		return nil
	},
}

var chatsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a chat",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat ID: %s", args[0])
		}

		var confirm bool
		if err := huh.NewConfirm().
			Title(fmt.Sprintf("Delete chat %d?", id)).
			Description("This will delete the chat and all messages.").
			Value(&confirm).
			Run(); err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewChatService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Chat deleted.")
		return nil
	},
}

var chatsSendCmd = &cobra.Command{
	Use:   "send <id>",
	Short: "Send a message to a chat",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat ID: %s", args[0])
		}

		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			return fmt.Errorf("--message is required")
		}

		client := mustClient()
		svc := api.NewChatService(client)
		if err := svc.SendMessage(cmdContext(), id, message); err != nil {
			return err
		}
		fmt.Println("Message sent. Response is being generated asynchronously.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatsCmd)

	chatsCmd.AddCommand(chatsListCmd)
	chatsListCmd.Flags().String("persona-id", "", "filter by persona ID")
	var page, perPage int
	addPaginationFlags(chatsListCmd, &page, &perPage)

	chatsCmd.AddCommand(chatsShowCmd)

	chatsCmd.AddCommand(chatsCreateCmd)
	chatsCreateCmd.Flags().String("prompt", "", "initial prompt")
	chatsCreateCmd.Flags().Int64("persona-id", 0, "persona ID")
	chatsCreateCmd.Flags().String("model", "", "model override")

	chatsCmd.AddCommand(chatsDeleteCmd)

	chatsCmd.AddCommand(chatsSendCmd)
	chatsSendCmd.Flags().String("message", "", "message content")
}
