package cmd

import (
	"context"

	"github.com/kuritka/k8gb-discovery/internal/common/depresolver"
	"github.com/kuritka/k8gb-discovery/internal/common/guard"
	"github.com/kuritka/k8gb-discovery/internal/common/runner"
	"github.com/kuritka/k8gb-discovery/internal/services/discovery"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "k8gb config discovery",
	Long:  `based on key: [namespace, cluster] returns k8gb configuration as marshalled json object 🍻🕺`,

	Run: func(cmd *cobra.Command, args []string) {
		settings, err := depresolver.New().MustResolveDiscovery()
		guard.FailOnError(err, "resolve input settings")
		disco := discovery.New(context.Background(), settings)
		runner.New(disco).MustRun()
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}
