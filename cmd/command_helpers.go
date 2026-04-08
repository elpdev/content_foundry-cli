package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func isInteractiveTerminal() bool {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return false
	}
	_ = tty.Close()
	return true
}

func addAutoConfirmFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("auto-confirm", false, "skip interactive confirmation prompt")
	cmd.Flags().Bool("yes", false, "alias for --auto-confirm")
}

func confirmDestructiveAction(cmd *cobra.Command, title, description string) (bool, error) {
	autoConfirm, _ := cmd.Flags().GetBool("auto-confirm")
	yes, _ := cmd.Flags().GetBool("yes")
	if autoConfirm || yes {
		return true, nil
	}

	if !isInteractiveTerminal() {
		return false, fmt.Errorf("interactive confirmation requires a TTY; rerun with --auto-confirm")
	}

	var confirm bool
	if err := huh.NewConfirm().
		Title(title).
		Description(description).
		Value(&confirm).
		Run(); err != nil {
		return false, err
	}

	return confirm, nil
}
