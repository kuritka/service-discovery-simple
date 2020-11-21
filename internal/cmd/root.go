// Package cmd implements one command
// - listen, listen for k8gb initial requests and returns disco 🕺🕺
package cmd

import (
	"fmt"
	"os"

	"github.com/kuritka/k8gb-discovery/internal/common"
	"github.com/kuritka/k8gb-discovery/internal/common/guard"

	"github.com/enescakir/emoji"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	// Verbose output
	Verbose bool
)

var rootCmd = &cobra.Command{
	Short: "k8gb discovery",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			guard.Message("No parameters included")
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n Not sure what to do %s? check out %s! %s\n", aurora.BrightGreen("next"), aurora.BrightBlue(common.HomeURL), emoji.BeachWithUmbrella)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

// Execute runs concrete command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
