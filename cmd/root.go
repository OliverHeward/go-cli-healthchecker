package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "A CLI tool for health checking endpoints",
	Long: `Health Checker is a CLI tool that pings multiple endpoints 
and reports their health status with response times.

Run 'healthcheck check' to perform health checks on configured endpoints.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
