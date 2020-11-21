// Package cmd implements one command
// - listen, listen for k8gb initial requests and returns disco 🕺🕺
package cmd

import (
	"os"

	"github.com/kuritka/k8gb-discovery/internal/common/log"

	"github.com/kuritka/k8gb-discovery/internal/common"
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
			log.Logger().Info("No parameters included")
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Logger().Info("Not sure what to do next? check out ", common.HomeURL)
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
