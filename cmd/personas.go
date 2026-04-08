package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

var personasCmd = &cobra.Command{
	Use:   "personas",
	Short: "Manage AI personas",
}

var personasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List personas",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		svc := api.NewPersonaService(client)

		personas, err := svc.List(cmdContext())
		if err != nil {
			return err
		}

		headers := []string{"ID", "Name", "Role", "Emoji", "Active"}
		rows := make([][]string, len(personas))
		for i, p := range personas {
			rows[i] = []string{
				fmt.Sprintf("%d", p.ID), p.Name, p.RoleTitle, p.AvatarEmoji,
				fmt.Sprintf("%t", p.Active),
			}
		}

		fmt.Print(formatter.FormatList(headers, rows, nil))
		return nil
	},
}

var personasShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show persona details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid persona ID: %s", args[0])
		}

		client := mustClient()
		svc := api.NewPersonaService(client)
		persona, err := svc.Get(cmdContext(), id)
		if err != nil {
			return err
		}

		prompt := persona.SystemPrompt
		if len(prompt) > 200 {
			prompt = prompt[:200] + "..."
		}

		fields := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", persona.ID)},
			{Key: "Name", Value: persona.Name},
			{Key: "Role", Value: persona.RoleTitle},
			{Key: "Description", Value: persona.Description},
			{Key: "Emoji", Value: persona.AvatarEmoji},
			{Key: "Model", Value: persona.ModelDisplayName()},
			{Key: "Active", Value: fmt.Sprintf("%t", persona.Active)},
			{Key: "System Prompt", Value: prompt},
			{Key: "Created", Value: persona.CreatedAt},
		}
		fmt.Print(formatter.FormatItem(fields))
		return nil
	},
}

var personasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a persona",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		role, _ := cmd.Flags().GetString("role")
		prompt, _ := cmd.Flags().GetString("prompt")

		if name == "" {
			if !isInteractiveTerminal() {
				return fmt.Errorf("--name is required")
			}
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Persona name").Value(&name),
					huh.NewInput().Title("Role title").Value(&role),
					huh.NewText().Title("System prompt").Value(&prompt),
				),
			).WithTheme(huh.ThemeCharm())
			if err := form.Run(); err != nil {
				return err
			}
		}

		fields := map[string]any{"name": name}
		if role != "" {
			fields["role_title"] = role
		}
		if prompt != "" {
			fields["system_prompt"] = prompt
		}

		client := mustClient()
		svc := api.NewPersonaService(client)
		persona, err := svc.Create(cmdContext(), fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", persona.ID)},
			{Key: "Name", Value: persona.Name},
			{Key: "Role", Value: persona.RoleTitle},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var personasUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a persona",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid persona ID: %s", args[0])
		}

		fields := map[string]any{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			fields["name"] = v
		}
		if cmd.Flags().Changed("role") {
			v, _ := cmd.Flags().GetString("role")
			fields["role_title"] = v
		}
		if cmd.Flags().Changed("prompt") {
			v, _ := cmd.Flags().GetString("prompt")
			fields["system_prompt"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			fields["description"] = v
		}
		if cmd.Flags().Changed("model") {
			v, _ := cmd.Flags().GetString("model")
			fields["model_id"] = v
		}

		if len(fields) == 0 {
			return fmt.Errorf("no fields to update (use --name, --role, --prompt, --description, --model)")
		}

		client := mustClient()
		svc := api.NewPersonaService(client)
		persona, err := svc.Update(cmdContext(), id, fields)
		if err != nil {
			return err
		}

		out := []output.Field{
			{Key: "ID", Value: fmt.Sprintf("%d", persona.ID)},
			{Key: "Name", Value: persona.Name},
			{Key: "Role", Value: persona.RoleTitle},
		}
		fmt.Print(formatter.FormatItem(out))
		return nil
	},
}

var personasDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a persona",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid persona ID: %s", args[0])
		}

		confirm, err := confirmDestructiveAction(cmd, fmt.Sprintf("Delete persona %d?", id), "This cannot be undone.")
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}

		client := mustClient()
		svc := api.NewPersonaService(client)
		if err := svc.Delete(cmdContext(), id); err != nil {
			return err
		}
		fmt.Println("Persona deleted.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(personasCmd)

	personasCmd.AddCommand(personasListCmd)
	personasCmd.AddCommand(personasShowCmd)

	personasCmd.AddCommand(personasCreateCmd)
	personasCreateCmd.Flags().String("name", "", "persona name")
	personasCreateCmd.Flags().String("role", "", "role title")
	personasCreateCmd.Flags().String("prompt", "", "system prompt")

	personasCmd.AddCommand(personasUpdateCmd)
	personasUpdateCmd.Flags().String("name", "", "new name")
	personasUpdateCmd.Flags().String("role", "", "new role title")
	personasUpdateCmd.Flags().String("prompt", "", "new system prompt")
	personasUpdateCmd.Flags().String("description", "", "new description")
	personasUpdateCmd.Flags().String("model", "", "AI model provider ID (e.g. claude-sonnet-4-5)")

	personasCmd.AddCommand(personasDeleteCmd)
	addAutoConfirmFlags(personasDeleteCmd)
}
