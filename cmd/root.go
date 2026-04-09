package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/leo/content-foundry-cli/internal/api"
	"github.com/leo/content-foundry-cli/internal/config"
	"github.com/leo/content-foundry-cli/internal/logging"
	"github.com/leo/content-foundry-cli/internal/output"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var (
	cfgPath   string
	outFormat string
	verbose   bool
	brandFlag string

	cfg       *config.Config
	apiClient *api.Client
	formatter output.Formatter
)

var rootCmd = &cobra.Command{
	Use:   "content_foundry",
	Short: "Content Foundry CLI & TUI",
	Long: lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2d95")).Bold(true).Render("CONTENT FOUNDRY") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#6c6c8a")).Render(" -- manage your content pipeline from the terminal"),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetVerbose(verbose)
		formatter = output.New(outFormat)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file path (default: ~/.config/content_foundry/config.toml)")
	rootCmd.PersistentFlags().StringVarP(&outFormat, "format", "f", "table", "output format: table, json, text")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&brandFlag, "brand", "", "brand ID or slug override")
}

// mustLoadConfig loads the config or exits with an error.
func mustLoadConfig() *config.Config {
	if cfg != nil {
		return cfg
	}
	var err error
	if cfgPath != "" {
		cfg, err = config.LoadFrom(cfgPath)
	} else {
		cfg, err = config.Load()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Render(
			"Error loading config: "+err.Error(),
		))
		fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#6c6c8a")).Render(
			"Run 'content_foundry auth login' to set up your credentials.",
		))
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Render(
			"Invalid config: "+err.Error(),
		))
		os.Exit(1)
	}
	return cfg
}

// mustClient returns an authenticated API client or exits.
func mustClient() *api.Client {
	if apiClient != nil {
		return apiClient
	}
	c := mustLoadConfig()
	apiClient = api.NewClient(c)

	// Override brand if --brand flag is set
	if brandFlag != "" {
		if id, err := strconv.ParseInt(brandFlag, 10, 64); err == nil {
			apiClient.BrandID = id
		} else {
			resolverClient := api.NewClient(c)
			resolverClient.BrandID = 0
			brand, err := api.NewBrandService(resolverClient).GetByRef(cmdContext(), brandFlag)
			if err != nil {
				fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Render(
					"Error resolving brand override: "+err.Error(),
				))
				os.Exit(1)
			}
			apiClient.BrandID = brand.ID
		}
	}

	return apiClient
}

// cmdContext returns a context for CLI commands.
func cmdContext() context.Context {
	return context.Background()
}

// addPaginationFlags adds --page and --per-page flags to a command.
func addPaginationFlags(cmd *cobra.Command, page, perPage *int) {
	cmd.Flags().IntVar(page, "page", 1, "page number")
	cmd.Flags().IntVar(perPage, "per-page", 20, "items per page (max 100)")
}
