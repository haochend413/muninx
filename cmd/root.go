package cmd

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "munina",
	Short: "munina",
	Long:  "Munina",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get state (can be nil if first run)
		s, err := state.LoadState(globalCfg.StateFilePath)
		if err != nil {
			// Use default state if load fails
			s = state.DefaultState()
		}

		// Initialize application with AppState
		globalApp = app.NewApp(globalDB, &s.App)

		// Initialize UI model with full state
		model := ui.NewModel(globalApp, globalCfg, s)
		globalModel = &model

		// Run Bubble Tea program
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	defer func() {
		if globalDB != nil {
			globalDB.Close()
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
}
