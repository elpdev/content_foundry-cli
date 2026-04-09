package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Set up API credentials",
	RunE:  runAuthLogin,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff2d95")).
		Bold(true)

	fmt.Println(titleStyle.Render("CONTENT FOUNDRY") + " -- connect your account")
	fmt.Println()

	var baseURL, clientID, secretKey string

	existing, _ := config.Load()
	if existing != nil {
		baseURL = existing.BaseURL
		clientID = existing.ClientID
	}

	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	theme := huh.ThemeCharm()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Base URL").
				Description("Your Content Foundry server URL").
				Value(&baseURL),

			huh.NewInput().
				Title("Client ID").
				Description("API key client_id").
				Value(&clientID),

			huh.NewInput().
				Title("Secret Key").
				Description("API key secret").
				EchoMode(huh.EchoModePassword).
				Value(&secretKey),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		return err
	}

	newCfg := &config.Config{
		BaseURL:   baseURL,
		ClientID:  clientID,
		SecretKey: secretKey,
	}

	if err := newCfg.Validate(); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	// Test credentials by listing brands
	fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#6c6c8a")).Render("Testing credentials... "))
	client := api.NewClient(newCfg)
	svc := api.NewBrandService(client)
	brands, err := svc.List(cmdContext(), api.BrandListParams{Page: 1, PerPage: 5})
	if err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Render("failed"))
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#44ff44")).Render("ok"))
	fmt.Println()
	fmt.Printf("  Found %d brand(s)\n", brands.Pagination.Count)

	// If there's exactly one brand, set it as default
	if len(brands.Items) == 1 {
		newCfg.DefaultBrandID = brands.Items[0].ID
		newCfg.DefaultBrandSlug = brands.Items[0].Slug
		fmt.Printf("  Default brand: %s (%s)\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#00fff2")).Bold(true).Render(brands.Items[0].Name),
			brands.Items[0].Slug,
		)
	} else if len(brands.Items) > 1 {
		fmt.Printf("  Run 'content_foundry brands use <id|slug>' to set your default brand.\n")
	}

	if err := newCfg.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not save config:", err)
	} else {
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#6c6c8a")).Render(
			"  Config saved to " + config.ConfigPath(),
		))
	}

	return nil
}
